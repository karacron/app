# Assistant Message Schema

Estructura de datos para enviar y recibir respuestas de un asistente multimodal.

Este documento define un contrato tipado para mensajes entre usuario, asistente, sistema y herramientas. Está pensado para asistentes que pueden responder con texto, audio, documentos, imágenes, vídeos, respuestas booleanas, solicitudes de permisos, multipreguntas, razonamiento resumido y estados de procesamiento.

---

## Objetivo

La estructura busca resolver estos casos:

- Enviar texto simple.
- Enviar respuestas multimedia.
- Enviar audios.
- Enviar documentos o cualquier tipo de archivo.
- Referenciar la pregunta original o documentos usados.
- Responder preguntas de sí o no.
- Solicitar permisos al usuario.
- Gestionar multipreguntas.
- Permitir respuestas del asistente con estados como `thinking` y `reasoning`.
- Mantener compatibilidad con TypeScript, Zod, WebSocket, SSE, bases de datos y frontend dinámico.

---

## Principio de diseño

Evitar estructuras demasiado simples como esta:

```ts
{
	type: 'text',
	data: 'Texto de ejemplo',
}
```

Es mejor usar un payload extensible:

```ts
{
	type: 'text',
	data: {
		text: 'Texto de ejemplo',
	},
}
```

Esto permite añadir campos futuros sin romper el contrato:

```ts
{
	type: 'text',
	data: {
		text: 'Texto de ejemplo',
		format: 'markdown',
		language: 'es',
	},
}
```

---

## Estructura base

```ts
export type AssistantMessageRole = "user" | "assistant" | "system" | "tool";

export type AssistantMessageType =
	| "text"
	| "audio"
	| "image"
	| "video"
	| "document"
	| "file"
	| "boolean"
	| "permission_request"
	| "permission_response"
	| "multi_question"
	| "multi_answer"
	| "thinking"
	| "reasoning"
	| "error"
	| "event";

export interface AssistantMessage<
	TType extends AssistantMessageType = AssistantMessageType,
	TData = unknown,
> {
	id: string;
	conversationId: string;
	role: AssistantMessageRole;
	type: TType;
	data: TData;
	references?: MessageReference[];
	metadata?: MessageMetadata;
	stream?: MessageStreamMetadata;
	createdAt: string;
}

export type NonChunkableAssistantMessageType =
	| "permission_request"
	| "multi_question"
	| "boolean"
	| "thinking";

export type ChunkableAssistantMessageType = Exclude<
	AssistantMessageType,
	NonChunkableAssistantMessageType
>;

export interface MessageStreamMetadata {
	mode: "chunked" | "atomic";
	isChunk?: boolean;
	chunkIndex?: number;
	chunkCount?: number;
	isFinal?: boolean;
	groupId?: string;
}
```

Regla de transporte:

- Las respuestas del asistente pueden enviarse en chunks para no esperar a que termine de escribir toda la salida.
- Los tipos `permission_request`, `multi_question`, `boolean` y `thinking` deben enviarse siempre completos, nunca fragmentados.
- Cuando un mensaje viaja fragmentado, el tipo lógico final no cambia; lo que cambia es el metadato `stream` y los eventos realtime del transporte.

---

## Referencias

Las referencias permiten vincular una respuesta con:

- La pregunta original.
- Otro mensaje.
- Un documento.
- Un archivo.
- Un chunk de contexto.
- Una llamada a herramienta.
- Una URL externa.
- Una solicitud de permisos.

```ts
export type ReferenceType =
	| "message"
	| "question"
	| "answer"
	| "document"
	| "file"
	| "chunk"
	| "tool_call"
	| "permission"
	| "external_url";

export interface MessageReference {
	id: string;
	type: ReferenceType;
	label?: string;
	source?: string;
	url?: string;
	metadata?: Record<string, unknown>;
}
```

Ejemplo:

```ts
const references: MessageReference[] = [
	{
		id: "msg_123",
		type: "question",
		label: "Pregunta original",
	},
	{
		id: "doc_456",
		type: "document",
		label: "Contrato PDF",
	},
];
```

---

## Metadata

```ts
export interface MessageMetadata {
	model?: string;
	provider?: string;
	locale?: string;
	visibility?: "public" | "private" | "internal";
	source?: "user_input" | "assistant_output" | "tool_output" | "system";
	latencyMs?: number;
	tokens?: {
		input?: number;
		output?: number;
		total?: number;
	};
	[key: string]: unknown;
}
```

Ejemplo:

```ts
const metadata: MessageMetadata = {
	model: "qwen3-4b-instruct",
	provider: "local",
	locale: "es",
	visibility: "public",
	source: "assistant_output",
	latencyMs: 320,
	tokens: {
		input: 120,
		output: 80,
		total: 200,
	},
};
```

---

# Tipos de mensajes

## 1. Texto

```ts
export interface TextMessageData {
	text: string;
	format?: "plain" | "markdown" | "html";
	language?: string;
}
```

Ejemplo:

```ts
const textMessage: AssistantMessage<"text", TextMessageData> = {
	id: "msg_001",
	conversationId: "conv_001",
	role: "assistant",
	type: "text",
	data: {
		text: "Texto de ejemplo",
		format: "markdown",
		language: "es",
	},
	createdAt: new Date().toISOString(),
};
```

JSON:

```json
{
	"id": "msg_001",
	"conversationId": "conv_001",
	"role": "assistant",
	"type": "text",
	"data": {
		"text": "Texto de ejemplo",
		"format": "markdown",
		"language": "es"
	},
	"createdAt": "2026-05-17T10:00:00.000Z"
}
```

---

## Streaming por chunks

Para respuestas largas del asistente, especialmente texto generado progresivamente, el servidor puede emitir varios fragmentos antes del mensaje final.

```ts
const chunkedTextMessage: AssistantMessage<"text", TextMessageData> = {
	id: "msg_200",
	conversationId: "conv_001",
	role: "assistant",
	type: "text",
	data: {
		text: "Primera parte de la respuesta...",
		format: "markdown",
		language: "es",
	},
	stream: {
		mode: "chunked",
		isChunk: true,
		chunkIndex: 0,
		chunkCount: 3,
		isFinal: false,
		groupId: "grp_001",
	},
	createdAt: new Date().toISOString(),
};
```

El último fragmento debe marcar `isFinal: true` o ir seguido por un evento `message.completed`.

Tipos atómicos que no deben fragmentarse:

- `permission_request`
- `multi_question`
- `boolean`
- `thinking`

---

## 2. Audio

```ts
export interface AudioMessageData {
	url: string;
	mimeType: string;
	durationMs?: number;
	transcript?: string;
	sizeBytes?: number;
}
```

Ejemplo:

```ts
const audioMessage: AssistantMessage<"audio", AudioMessageData> = {
	id: "msg_002",
	conversationId: "conv_001",
	role: "assistant",
	type: "audio",
	data: {
		url: "https://cdn.example.com/audio/response.mp3",
		mimeType: "audio/mpeg",
		durationMs: 12000,
		transcript: "Esta es la transcripción del audio.",
		sizeBytes: 500000,
	},
	createdAt: new Date().toISOString(),
};
```

---

## 3. Imagen

```ts
export interface ImageMessageData {
	url: string;
	mimeType: string;
	alt?: string;
	width?: number;
	height?: number;
	sizeBytes?: number;
}
```

Ejemplo:

```ts
const imageMessage: AssistantMessage<"image", ImageMessageData> = {
	id: "msg_003",
	conversationId: "conv_001",
	role: "assistant",
	type: "image",
	data: {
		url: "https://cdn.example.com/images/result.png",
		mimeType: "image/png",
		alt: "Imagen generada por el asistente",
		width: 1024,
		height: 1024,
	},
	createdAt: new Date().toISOString(),
};
```

---

## 4. Vídeo

```ts
export interface VideoMessageData {
	url: string;
	mimeType: string;
	durationMs?: number;
	thumbnailUrl?: string;
	transcript?: string;
	sizeBytes?: number;
}
```

Ejemplo:

```ts
const videoMessage: AssistantMessage<"video", VideoMessageData> = {
	id: "msg_004",
	conversationId: "conv_001",
	role: "assistant",
	type: "video",
	data: {
		url: "https://cdn.example.com/video/demo.mp4",
		mimeType: "video/mp4",
		durationMs: 60000,
		thumbnailUrl: "https://cdn.example.com/video/demo-thumb.jpg",
	},
	createdAt: new Date().toISOString(),
};
```

---

## 5. Documento o archivo

```ts
export interface DocumentMessageData {
	fileId: string;
	name: string;
	url: string;
	mimeType: string;
	extension?: string;
	sizeBytes?: number;
	description?: string;
}
```

Ejemplo:

```ts
const documentMessage: AssistantMessage<"document", DocumentMessageData> = {
	id: "msg_005",
	conversationId: "conv_001",
	role: "assistant",
	type: "document",
	data: {
		fileId: "file_001",
		name: "informe.pdf",
		url: "https://cdn.example.com/files/informe.pdf",
		mimeType: "application/pdf",
		extension: "pdf",
		sizeBytes: 204800,
		description: "Informe generado por el asistente",
	},
	createdAt: new Date().toISOString(),
};
```

---

## 6. Respuesta sí/no

```ts
export interface BooleanMessageData {
	value: boolean;
	label?: string;
	reason?: string;
	confidence?: number;
}
```

Ejemplo:

```ts
const booleanMessage: AssistantMessage<"boolean", BooleanMessageData> = {
	id: "msg_006",
	conversationId: "conv_001",
	role: "assistant",
	type: "boolean",
	data: {
		value: true,
		label: "Sí",
		reason: "El usuario tiene permisos suficientes.",
		confidence: 0.92,
	},
	createdAt: new Date().toISOString(),
};
```

---

## 7. Solicitud de permisos

Sirve cuando el asistente necesita autorización explícita del usuario.

Ejemplos:

- Leer archivo.
- Escribir archivo.
- Ejecutar comando.
- Enviar email.
- Acceder al calendario.
- Usar micrófono.
- Usar cámara.
- Acceder a ubicación.
- Controlar navegador.
- Llamar a una API externa.

```ts
export type PermissionAction =
	| "read_file"
	| "write_file"
	| "delete_file"
	| "execute_command"
	| "send_email"
	| "read_calendar"
	| "write_calendar"
	| "access_microphone"
	| "access_camera"
	| "access_location"
	| "external_api_call"
	| "browser_control";

export interface PermissionRequestMessageData {
	permissionId: string;
	action: PermissionAction;
	title: string;
	description: string;
	required: boolean;
	scope?: string;
	expiresAt?: string;
	options?: PermissionOption[];
}

export interface PermissionOption {
	id: string;
	label: string;
	value: "allow" | "deny" | "allow_once" | "always_allow";
}
```

Ejemplo:

```ts
const permissionRequestMessage: AssistantMessage<
	"permission_request",
	PermissionRequestMessageData
> = {
	id: "msg_007",
	conversationId: "conv_001",
	role: "assistant",
	type: "permission_request",
	data: {
		permissionId: "perm_001",
		action: "execute_command",
		title: "Permiso para ejecutar comando",
		description:
			"Necesito ejecutar pnpm install para instalar las dependencias del proyecto.",
		required: true,
		scope: "current_project",
		options: [
			{
				id: "opt_1",
				label: "Permitir una vez",
				value: "allow_once",
			},
			{
				id: "opt_2",
				label: "Denegar",
				value: "deny",
			},
		],
	},
	createdAt: new Date().toISOString(),
};
```

---

## 8. Respuesta a permisos

```ts
export interface PermissionResponseMessageData {
	permissionId: string;
	value: "allow" | "deny" | "allow_once" | "always_allow";
	reason?: string;
}
```

Ejemplo:

```ts
const permissionResponseMessage: AssistantMessage<
	"permission_response",
	PermissionResponseMessageData
> = {
	id: "msg_008",
	conversationId: "conv_001",
	role: "user",
	type: "permission_response",
	data: {
		permissionId: "perm_001",
		value: "allow_once",
	},
	createdAt: new Date().toISOString(),
};
```

---

## 9. Multipreguntas

```ts
export interface MultiQuestionMessageData {
	questions: QuestionItem[];
	allowPartialAnswers?: boolean;
}

export interface QuestionItem {
	id: string;
	text: string;
	type:
		| "text"
		| "boolean"
		| "single_choice"
		| "multiple_choice"
		| "file"
		| "number"
		| "date";
	required?: boolean;
	options?: QuestionOption[];
}

export interface QuestionOption {
	id: string;
	label: string;
	value: string | number | boolean;
}
```

Ejemplo:

```ts
const multiQuestionMessage: AssistantMessage<
	"multi_question",
	MultiQuestionMessageData
> = {
	id: "msg_009",
	conversationId: "conv_001",
	role: "assistant",
	type: "multi_question",
	data: {
		allowPartialAnswers: true,
		questions: [
			{
				id: "q_001",
				text: "¿Quieres que el asistente ejecute comandos?",
				type: "boolean",
				required: true,
			},
			{
				id: "q_002",
				text: "¿Qué modelo local quieres usar?",
				type: "single_choice",
				required: true,
				options: [
					{
						id: "opt_qwen",
						label: "Qwen 3",
						value: "qwen3",
					},
					{
						id: "opt_llama",
						label: "Llama",
						value: "llama",
					},
				],
			},
		],
	},
	createdAt: new Date().toISOString(),
};
```

---

## 10. Respuesta a multipreguntas

```ts
export interface MultiAnswerMessageData {
	answers: AnswerItem[];
}

export interface AnswerItem {
	questionId: string;
	value: string | number | boolean | string[] | FileAnswerData | null;
}

export interface FileAnswerData {
	fileId: string;
	name: string;
	url: string;
	mimeType: string;
}
```

Ejemplo:

```ts
const multiAnswerMessage: AssistantMessage<
	"multi_answer",
	MultiAnswerMessageData
> = {
	id: "msg_010",
	conversationId: "conv_001",
	role: "user",
	type: "multi_answer",
	data: {
		answers: [
			{
				questionId: "q_001",
				value: true,
			},
			{
				questionId: "q_002",
				value: "qwen3",
			},
		],
	},
	createdAt: new Date().toISOString(),
};
```

---

## 11. Thinking

Representa que el asistente está procesando una tarea.

No debe usarse para exponer razonamiento interno sensible. Sirve para UI, progreso, streaming y trazabilidad visible.

```ts
export interface ThinkingMessageData {
	status: "started" | "in_progress" | "completed";
	text?: string;
	step?: string;
	progress?: number;
}
```

Ejemplo:

```ts
const thinkingMessage: AssistantMessage<"thinking", ThinkingMessageData> = {
	id: "msg_011",
	conversationId: "conv_001",
	role: "assistant",
	type: "thinking",
	data: {
		status: "in_progress",
		text: "Analizando documentos...",
		step: "document_analysis",
		progress: 45,
	},
	metadata: {
		visibility: "public",
	},
	createdAt: new Date().toISOString(),
};
```

---

## 12. Reasoning

Representa una explicación resumida del proceso de decisión.

Recomendación: guardar reasoning como resumen técnico o trazabilidad, no como cadena completa de pensamiento privada.

```ts
export interface ReasoningMessageData {
	summary: string;
	steps?: ReasoningStep[];
	confidence?: number;
}

export interface ReasoningStep {
	id: string;
	title: string;
	description: string;
	status: "pending" | "completed" | "failed";
}
```

Ejemplo:

```ts
const reasoningMessage: AssistantMessage<"reasoning", ReasoningMessageData> = {
	id: "msg_012",
	conversationId: "conv_001",
	role: "assistant",
	type: "reasoning",
	data: {
		summary:
			"He comparado las opciones disponibles y he elegido la más segura.",
		confidence: 0.87,
		steps: [
			{
				id: "step_001",
				title: "Validar permisos",
				description:
					"Se comprobó que la acción requiere autorización explícita.",
				status: "completed",
			},
			{
				id: "step_002",
				title: "Preparar respuesta",
				description:
					"Se generó una respuesta segura y compatible con el flujo.",
				status: "completed",
			},
		],
	},
	metadata: {
		visibility: "public",
	},
	createdAt: new Date().toISOString(),
};
```

---

## 13. Error

```ts
export interface ErrorMessageData {
	code: string;
	message: string;
	details?: Record<string, unknown>;
	retryable?: boolean;
}
```

Ejemplo:

```ts
const errorMessage: AssistantMessage<"error", ErrorMessageData> = {
	id: "msg_013",
	conversationId: "conv_001",
	role: "assistant",
	type: "error",
	data: {
		code: "TOOL_EXECUTION_FAILED",
		message: "No se pudo ejecutar la herramienta solicitada.",
		retryable: true,
	},
	createdAt: new Date().toISOString(),
};
```

---

## 14. Eventos

Útil para registrar acciones internas o notificar cambios al frontend.

```ts
export interface EventMessageData {
	name: string;
	payload?: Record<string, unknown>;
}
```

Ejemplo:

```ts
const eventMessage: AssistantMessage<"event", EventMessageData> = {
	id: "msg_014",
	conversationId: "conv_001",
	role: "system",
	type: "event",
	data: {
		name: "conversation.updated",
		payload: {
			messageCount: 14,
		},
	},
	createdAt: new Date().toISOString(),
};
```

---

# Unión tipada completa

```ts
export type AnyAssistantMessage =
	| AssistantMessage<"text", TextMessageData>
	| AssistantMessage<"audio", AudioMessageData>
	| AssistantMessage<"image", ImageMessageData>
	| AssistantMessage<"video", VideoMessageData>
	| AssistantMessage<"document", DocumentMessageData>
	| AssistantMessage<"file", DocumentMessageData>
	| AssistantMessage<"boolean", BooleanMessageData>
	| AssistantMessage<"permission_request", PermissionRequestMessageData>
	| AssistantMessage<"permission_response", PermissionResponseMessageData>
	| AssistantMessage<"multi_question", MultiQuestionMessageData>
	| AssistantMessage<"multi_answer", MultiAnswerMessageData>
	| AssistantMessage<"thinking", ThinkingMessageData>
	| AssistantMessage<"reasoning", ReasoningMessageData>
	| AssistantMessage<"error", ErrorMessageData>
	| AssistantMessage<"event", EventMessageData>;
```

---

# Ejemplo real: respuesta multimedia con referencias

```ts
const response: AnyAssistantMessage[] = [
	{
		id: "msg_100",
		conversationId: "conv_001",
		role: "assistant",
		type: "text",
		data: {
			text: "He generado el informe y también te dejo una explicación en audio.",
			format: "plain",
		},
		references: [
			{
				id: "msg_099",
				type: "question",
				label: "Pregunta original",
			},
		],
		createdAt: new Date().toISOString(),
	},
	{
		id: "msg_101",
		conversationId: "conv_001",
		role: "assistant",
		type: "document",
		data: {
			fileId: "file_123",
			name: "informe-final.pdf",
			url: "https://cdn.example.com/informe-final.pdf",
			mimeType: "application/pdf",
			extension: "pdf",
		},
		references: [
			{
				id: "msg_099",
				type: "question",
			},
		],
		createdAt: new Date().toISOString(),
	},
	{
		id: "msg_102",
		conversationId: "conv_001",
		role: "assistant",
		type: "audio",
		data: {
			url: "https://cdn.example.com/audio/resumen.mp3",
			mimeType: "audio/mpeg",
			transcript: "Resumen del informe generado.",
		},
		references: [
			{
				id: "file_123",
				type: "document",
				label: "Informe final",
			},
		],
		createdAt: new Date().toISOString(),
	},
];
```

---

# Ejemplo JSON completo

```json
{
	"id": "msg_001",
	"conversationId": "conv_001",
	"role": "assistant",
	"type": "text",
	"data": {
		"text": "Texto de ejemplo",
		"format": "plain"
	},
	"references": [
		{
			"id": "msg_000",
			"type": "question",
			"label": "Pregunta original"
		}
	],
	"metadata": {
		"model": "qwen3-4b-instruct",
		"provider": "local",
		"visibility": "public"
	},
	"createdAt": "2026-05-17T10:00:00.000Z"
}
```

---

# Validación con Zod

Recomendado si quieres blindar entradas y salidas.

```ts
import { z } from "zod";

export const messageReferenceSchema = z.object({
	id: z.string(),
	type: z.enum([
		"message",
		"question",
		"answer",
		"document",
		"file",
		"chunk",
		"tool_call",
		"permission",
		"external_url",
	]),
	label: z.string().optional(),
	source: z.string().optional(),
	url: z.string().url().optional(),
	metadata: z.record(z.string(), z.unknown()).optional(),
});

export const textMessageDataSchema = z.object({
	text: z.string(),
	format: z.enum(["plain", "markdown", "html"]).optional(),
	language: z.string().optional(),
});

export const assistantMessageBaseSchema = z.object({
	id: z.string(),
	conversationId: z.string(),
	role: z.enum(["user", "assistant", "system", "tool"]),
	type: z.string(),
	data: z.unknown(),
	references: z.array(messageReferenceSchema).optional(),
	metadata: z.record(z.string(), z.unknown()).optional(),
	createdAt: z.string().datetime(),
});

export const textMessageSchema = assistantMessageBaseSchema.extend({
	type: z.literal("text"),
	data: textMessageDataSchema,
});
```

---

# Renderizado en frontend

Ejemplo simple con React:

```tsx
export function MessageRenderer({ message }: { message: AnyAssistantMessage }) {
	switch (message.type) {
		case "text":
			return <p>{message.data.text}</p>;

		case "audio":
			return <audio controls src={message.data.url} />;

		case "image":
			return <img src={message.data.url} alt={message.data.alt ?? ""} />;

		case "video":
			return <video controls src={message.data.url} />;

		case "document":
		case "file":
			return (
				<a href={message.data.url} target="_blank" rel="noreferrer">
					{message.data.name}
				</a>
			);

		case "boolean":
			return <span>{message.data.value ? "Sí" : "No"}</span>;

		case "permission_request":
			return <PermissionRequest message={message} />;

		case "multi_question":
			return <MultiQuestionForm message={message} />;

		case "thinking":
			return <span>{message.data.text ?? "Procesando..."}</span>;

		case "reasoning":
			return <ReasoningSummary message={message} />;

		case "error":
			return <span>{message.data.message}</span>;

		default:
			return null;
	}
}
```

---

# Recomendación para base de datos

Tabla `messages`:

```sql
CREATE TABLE messages (
	id TEXT PRIMARY KEY,
	conversation_id TEXT NOT NULL,
	role TEXT NOT NULL,
	type TEXT NOT NULL,
	data JSONB NOT NULL,
	references JSONB,
	metadata JSONB,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Índices recomendados:

```sql
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_type ON messages(type);
CREATE INDEX idx_messages_created_at ON messages(created_at);
```

Para PostgreSQL:

```sql
CREATE INDEX idx_messages_data_gin ON messages USING GIN (data);
CREATE INDEX idx_messages_metadata_gin ON messages USING GIN (metadata);
```

---

# Recomendación para streaming

Para SSE o WebSocket, puedes emitir eventos con el mismo contrato.

Ejemplo SSE:

```ts
event: message.created
data: {
	"id": "msg_001",
	"conversationId": "conv_001",
	"role": "assistant",
	"type": "text",
	"data": {
		"text": "Hola"
	},
	"createdAt": "2026-05-17T10:00:00.000Z"
}
```

Eventos útiles:

```ts
export type AssistantRealtimeEvent =
	| "message.created"
	| "message.chunk"
	| "message.updated"
	| "message.completed"
	| "message.failed"
	| "permission.requested"
	| "permission.resolved"
	| "tool.started"
	| "tool.completed"
	| "tool.failed";
```

Ejemplo de secuencia chunked para una respuesta del asistente:

```ts
event: message.created
data: {
	"id": "msg_201",
	"conversationId": "conv_001",
	"role": "assistant",
	"type": "text",
	"data": { "text": "Hola" },
	"stream": {
		"mode": "chunked",
		"isChunk": true,
		"chunkIndex": 0,
		"chunkCount": 2,
		"groupId": "grp_002"
	},
	"createdAt": "2026-05-17T10:00:00.000Z"
}

event: message.chunk
data: {
	"id": "msg_201",
	"conversationId": "conv_001",
	"role": "assistant",
	"type": "text",
	"data": { "text": ", mundo" },
	"stream": {
		"mode": "chunked",
		"isChunk": true,
		"chunkIndex": 1,
		"chunkCount": 2,
		"isFinal": true,
		"groupId": "grp_002"
	},
	"createdAt": "2026-05-17T10:00:01.000Z"
}

event: message.completed
data: {
	"id": "msg_201",
	"conversationId": "conv_001",
	"role": "assistant",
	"type": "text",
	"data": { "text": "Hola, mundo" },
	"stream": {
		"mode": "chunked",
		"isFinal": true,
		"groupId": "grp_002"
	},
	"createdAt": "2026-05-17T10:00:01.100Z"
}
```

Para `permission_request`, `multi_question`, `boolean` y `thinking`, el transporte debe usar `mode: "atomic"` y emitir el mensaje completo en un solo evento.

---

# Implementación de conexión en tiempo real

La conexión en tiempo real debe construirse sobre un socket persistente y un envelope estable para todos los mensajes.

Contrato recomendado del cliente hacia el servidor:

```ts
type ClientToServerEvent = {
	event: "message" | "chat" | "permission_response" | "multi_answer";
	data: unknown;
};
```

Contrato recomendado del servidor hacia el cliente:

```ts
type ServerToClientEventName =
	| "connected"
	| "message.created"
	| "message.chunk"
	| "message.completed"
	| "message.failed"
	| "error";
```

Flujo recomendado:

1. El cliente abre la conexión WebSocket a `/ws`.
2. El servidor responde con `connected` y un `clientId`.
3. El cliente envía un envelope con `event` y `data`.
4. El servidor valida el payload y decide si la respuesta será `atomic` o `chunked`.
5. Si la respuesta es larga y fragmentable, el servidor emite `message.created`, luego uno o más `message.chunk`, y cierra con `message.completed`.
6. Si el tipo es `permission_request`, `multi_question`, `boolean` o `thinking`, el servidor emite un único mensaje completo con `stream.mode = "atomic"`.

Recomendaciones prácticas:

- Usa `ws://127.0.0.1:3000/ws` en pruebas locales si tu herramienta cliente falla con `localhost` por resolución IPv4/IPv6.
- Mantén el socket abierto durante toda la conversación, no abras una conexión por cada mensaje.
- Usa `groupId` para reconstruir una misma respuesta fragmentada en frontend.
- Emite `message.failed` cuando la generación falle después de haber empezado el flujo.
- No uses chunking para payloads semánticamente indivisibles.

## Ejemplo de implementación en Go

Ejemplo simplificado con Fiber y `fasthttp/websocket`, alineado con este contrato:

```go
package infrastructure

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type incomingEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type textData struct {
	Text   string `json:"text"`
	Format string `json:"format,omitempty"`
}

type streamMeta struct {
	Mode       string `json:"mode"`
	IsChunk    bool   `json:"isChunk,omitempty"`
	ChunkIndex int    `json:"chunkIndex,omitempty"`
	ChunkCount int    `json:"chunkCount,omitempty"`
	IsFinal    bool   `json:"isFinal,omitempty"`
	GroupID    string `json:"groupId,omitempty"`
}

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true },
}

func registerRealtime(app *fiber.App) {
	app.Get("/ws", func(c fiber.Ctx) error {
		return upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
			clientID := uuid.New().String()
			send(conn, map[string]any{
				"event": "connected",
				"data": map[string]any{"clientId": clientID},
			})

			for {
				_, raw, err := conn.ReadMessage()
				if err != nil {
					return
				}

				var incoming incomingEvent
				if err := json.Unmarshal(raw, &incoming); err != nil {
					send(conn, map[string]any{
						"event": "error",
						"data": map[string]any{"message": "JSON inválido"},
					})
					continue
				}

				switch incoming.Event {
				case "chat":
					streamAssistantText(conn, "Esta es una respuesta larga del asistente que llega en tiempo real.")
				case "message":
					send(conn, map[string]any{
						"event": "message.created",
						"data": map[string]any{
							"id":   uuid.New().String(),
							"role": "assistant",
							"type": "thinking",
							"data": map[string]any{"status": "in_progress", "text": "Procesando..."},
							"stream": streamMeta{Mode: "atomic"},
						},
					})
				}
			}
		})
	})
}

func streamAssistantText(conn *websocket.Conn, content string) {
	messageID := uuid.New().String()
	groupID := uuid.New().String()
	chunks := chunkString(content, 24)

	send(conn, map[string]any{
		"event": "message.created",
		"data": map[string]any{
			"id":   messageID,
			"role": "assistant",
			"type": "text",
			"data": textData{Text: chunks[0], Format: "markdown"},
			"stream": streamMeta{
				Mode:       "chunked",
				IsChunk:    true,
				ChunkIndex: 0,
				ChunkCount: len(chunks),
				GroupID:    groupID,
			},
		},
	})

	for i := 1; i < len(chunks); i++ {
		time.Sleep(80 * time.Millisecond)
		send(conn, map[string]any{
			"event": "message.chunk",
			"data": map[string]any{
				"id":   messageID,
				"role": "assistant",
				"type": "text",
				"data": textData{Text: chunks[i], Format: "markdown"},
				"stream": streamMeta{
					Mode:       "chunked",
					IsChunk:    true,
					ChunkIndex: i,
					ChunkCount: len(chunks),
					IsFinal:    i == len(chunks)-1,
					GroupID:    groupID,
				},
			},
		})
	}

	send(conn, map[string]any{
		"event": "message.completed",
		"data": map[string]any{
			"id":   messageID,
			"role": "assistant",
			"type": "text",
			"data": textData{Text: strings.Join(chunks, ""), Format: "markdown"},
			"stream": streamMeta{
				Mode:    "chunked",
				IsFinal: true,
				GroupID: groupID,
			},
		},
	})
}

func chunkString(text string, size int) []string {
	runes := []rune(text)
	if len(runes) == 0 {
		return []string{""}
	}
	chunks := make([]string, 0, (len(runes)+size-1)/size)
	for start := 0; start < len(runes); start += size {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}

func send(conn *websocket.Conn, payload any) {
	b, _ := json.Marshal(payload)
	_ = conn.WriteMessage(websocket.TextMessage, b)
}
```

## Ejemplo equivalente en NestJS

En NestJS la idea es la misma, pero usando un Gateway y `client.emit(...)` por evento:

```ts
import {
	ConnectedSocket,
	MessageBody,
	OnGatewayConnection,
	SubscribeMessage,
	WebSocketGateway,
	WebSocketServer,
} from "@nestjs/websockets";
import { Server, Socket } from "socket.io";
import { randomUUID } from "node:crypto";

@WebSocketGateway({ cors: { origin: "*" } })
export class AssistantGateway implements OnGatewayConnection {
	@WebSocketServer()
	server: Server;

	handleConnection(client: Socket) {
		client.emit("connected", {
			clientId: client.id,
		});
	}

	@SubscribeMessage("chat")
	async handleChat(
		@ConnectedSocket() client: Socket,
		@MessageBody() payload: { message: string },
	) {
		const messageId = randomUUID();
		const groupId = randomUUID();
		const text =
			"Esta es una respuesta larga del asistente enviada en tiempo real.";
		const chunks = this.chunkString(text, 24);

		client.emit("message.created", {
			id: messageId,
			role: "assistant",
			type: "text",
			data: { text: chunks[0], format: "markdown" },
			stream: {
				mode: "chunked",
				isChunk: true,
				chunkIndex: 0,
				chunkCount: chunks.length,
				groupId,
			},
		});

		for (let i = 1; i < chunks.length; i++) {
			await new Promise((resolve) => setTimeout(resolve, 80));
			client.emit("message.chunk", {
				id: messageId,
				role: "assistant",
				type: "text",
				data: { text: chunks[i], format: "markdown" },
				stream: {
					mode: "chunked",
					isChunk: true,
					chunkIndex: i,
					chunkCount: chunks.length,
					isFinal: i === chunks.length - 1,
					groupId,
				},
			});
		}

		client.emit("message.completed", {
			id: messageId,
			role: "assistant",
			type: "text",
			data: { text, format: "markdown" },
			stream: {
				mode: "chunked",
				isFinal: true,
				groupId,
			},
		});
	}

	@SubscribeMessage("request_permission")
	handlePermission(@ConnectedSocket() client: Socket) {
		client.emit("message.created", {
			id: randomUUID(),
			role: "assistant",
			type: "permission_request",
			data: {
				permissionId: randomUUID(),
				action: "execute_command",
				title: "Permiso para ejecutar comando",
				description: "Necesito permiso para continuar.",
				required: true,
			},
			stream: {
				mode: "atomic",
			},
		});
	}

	private chunkString(text: string, size: number): string[] {
		if (!text) return [""];
		const chars = Array.from(text);
		const chunks: string[] = [];
		for (let start = 0; start < chars.length; start += size) {
			chunks.push(chars.slice(start, start + size).join(""));
		}
		return chunks;
	}
}
```

## Cómo debe reconstruirlo el frontend

El cliente debe mantener un buffer por `groupId`:

1. `message.created`: crea el mensaje visible y guarda el primer fragmento.
2. `message.chunk`: concatena el texto al buffer existente.
3. `message.completed`: sustituye el contenido parcial por el final consolidado o marca el mensaje como terminado.
4. `message.failed`: marca error y conserva lo ya recibido si tiene sentido para la UX.

Para tipos atómicos, el frontend no debe esperar chunks adicionales.

---

# Buenas prácticas

## 1. Usar `type` como discriminador

El frontend y backend deben decidir cómo interpretar el payload usando `type`.

```ts
if (message.type === "text") {
	console.log(message.data.text);
}
```

---

## 2. Mantener `data` siempre como objeto

Evita esto:

```ts
data: "Texto plano";
```

Mejor:

```ts
data: {
	text: 'Texto plano',
}
```

---

## 3. Separar `thinking` de `reasoning`

`thinking` debe servir para estado de UI:

```ts
{
	type: 'thinking',
	data: {
		status: 'in_progress',
		text: 'Analizando archivos...',
	}
}
```

`reasoning` debe servir para explicación resumida:

```ts
{
	type: 'reasoning',
	data: {
		summary: 'Se eligió esta opción porque reduce riesgo y mantiene compatibilidad.',
	}
}
```

## 3.1. Fragmentar solo cuando aporte latencia percibida

Usa chunks cuando el usuario se beneficie de ver texto incrementalmente.

No fragmentes:

- `permission_request`
- `multi_question`
- `boolean`
- `thinking`

Esos tipos deben llegar completos para simplificar la UI y evitar estados intermedios ambiguos.

---

## 4. No exponer razonamiento interno sensible

No guardes cadenas privadas completas de pensamiento como respuesta pública.

Correcto:

```ts
{
	type: 'reasoning',
	data: {
		summary: 'Se revisaron tres alternativas y se eligió la más segura.',
	}
}
```

Incorrecto:

```ts
{
	type: 'reasoning',
	data: {
		internalChainOfThought: '...'
	}
}
```

---

## 5. Usar referencias para trazabilidad

Cuando una respuesta dependa de otra pregunta, documento o herramienta, usa `references`.

```ts
references: [
	{
		id: "msg_001",
		type: "question",
		label: "Pregunta original",
	},
];
```

---

## 6. Versionar el contrato

En producción conviene añadir versión al envelope o metadata.

```ts
metadata: {
	schemaVersion: '1.0.0',
}
```

---

# Estructura recomendada de carpetas

```txt
src/
├── messages/
│   ├── domain/
│   │   ├── assistant-message.ts
│   │   ├── message-reference.ts
│   │   ├── message-metadata.ts
│   │   └── message-types.ts
│   ├── schemas/
│   │   ├── assistant-message.schema.ts
│   │   ├── text-message.schema.ts
│   │   ├── media-message.schema.ts
│   │   ├── permission-message.schema.ts
│   │   └── multi-question.schema.ts
│   ├── application/
│   │   ├── create-message.use-case.ts
│   │   ├── validate-message.use-case.ts
│   │   └── render-message.use-case.ts
│   └── infrastructure/
│       ├── message.repository.ts
│       └── message.mapper.ts
```

---

# Conclusión

La estructura recomendada es:

```ts
AssistantMessage<TType, TData>;
```

Con:

- `type` como discriminador.
- `data` como payload tipado.
- `references` para trazabilidad.
- `metadata` para información técnica.
- `createdAt` para orden temporal.
- `role` para diferenciar usuario, asistente, sistema y herramientas.

Esta arquitectura permite crecer hacia un asistente local-first, multimodal, seguro y preparado para producción.
