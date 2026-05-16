INSERT OR IGNORE INTO installation (reference, version, platform, is_onboarded)
VALUES ('default', '0.1.0', 'windows', 0);

INSERT OR IGNORE INTO assistant_personalities (reference, domain, skill_reference, is_default, is_active)
VALUES
	('general', 'general', NULL, 1, 1),
	('backend-engineer', 'engineering', 'backend-engineer', 0, 1),
	('frontend-engineer', 'engineering', 'frontend-engineer', 0, 1),
	('fullstack-engineer', 'engineering', 'fullstack-engineer', 0, 1),
	('data-engineer', 'engineering', 'data-engineer', 0, 1),
	('ml-ai-engineer', 'engineering', 'ml-ai-engineer', 0, 1),
	('mobile-engineer', 'engineering', 'mobile-engineer', 0, 1),
	('devops-engineer', 'engineering', 'devops-engineer', 0, 1),
	('sre-engineer', 'engineering', 'sre-engineer', 0, 1),
	('platform-engineer', 'engineering', 'platform-engineer', 0, 1),
	('software-architect', 'engineering', 'software-architect', 0, 1),
	('qa-automation-engineer', 'engineering', 'qa-automation-engineer', 0, 1),
	('security-engineer', 'engineering', 'security-engineer', 0, 1),
	('creative-director', 'creative', 'creative-director', 0, 1),
	('designer', 'creative', 'designer', 0, 1),
	('manager', 'management', 'manager', 0, 1),
	('administrator', 'management', 'administrator', 0, 1),
	('mathematician', 'science', 'mathematician', 0, 1),
	('physicist', 'science', 'physicist', 0, 1),
	('biologist', 'science', 'biologist', 0, 1),
	('chemist', 'science', 'chemist', 0, 1),
	('physician', 'science', 'physician', 0, 1),
	('aeronautics-engineer', 'science', 'aeronautics-engineer', 0, 1),
	('astronomer', 'science', 'astronomer', 0, 1),
	('musician', 'arts', 'musician', 0, 1),
	('singer', 'arts', 'singer', 0, 1),
	('philosopher', 'humanities', 'philosopher', 0, 1),
	('marketing-sales', 'business', 'marketing-sales', 0, 1),
	('seo-sem-specialist', 'business', 'seo-sem-specialist', 0, 1),
	('productivity-workflows', 'tools', 'productivity-workflows', 0, 1);

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'en', 'General', 'Balanced assistant for day-to-day tasks and general support.'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'en',
	CASE reference
		WHEN 'backend-engineer' THEN 'Backend Engineer'
		WHEN 'frontend-engineer' THEN 'Frontend Engineer'
		WHEN 'fullstack-engineer' THEN 'Fullstack Engineer'
		WHEN 'data-engineer' THEN 'Data Engineer'
		WHEN 'ml-ai-engineer' THEN 'ML/AI Engineer'
		WHEN 'mobile-engineer' THEN 'Mobile Engineer'
		WHEN 'devops-engineer' THEN 'DevOps Engineer'
		WHEN 'sre-engineer' THEN 'SRE Engineer'
		WHEN 'platform-engineer' THEN 'Platform Engineer'
		WHEN 'software-architect' THEN 'Software Architect'
		WHEN 'qa-automation-engineer' THEN 'QA Automation Engineer'
		WHEN 'security-engineer' THEN 'Security Engineer'
		WHEN 'creative-director' THEN 'Creative Director'
		WHEN 'designer' THEN 'Designer'
		WHEN 'manager' THEN 'Manager'
		WHEN 'administrator' THEN 'Administrator'
		WHEN 'mathematician' THEN 'Mathematician'
		WHEN 'physicist' THEN 'Physicist'
		WHEN 'biologist' THEN 'Biologist'
		WHEN 'chemist' THEN 'Chemist'
		WHEN 'physician' THEN 'Physician'
		WHEN 'aeronautics-engineer' THEN 'Aeronautics Engineer'
		WHEN 'astronomer' THEN 'Astronomer'
		WHEN 'musician' THEN 'Musician'
		WHEN 'singer' THEN 'Singer'
		WHEN 'philosopher' THEN 'Philosopher'
		WHEN 'marketing-sales' THEN 'Marketing & Sales'
		WHEN 'seo-sem-specialist' THEN 'SEO/SEM Specialist'
		WHEN 'productivity-workflows' THEN 'Productivity Workflows'
	END,
	CASE reference
		WHEN 'backend-engineer' THEN 'Design and implement backend services with robust APIs, domain boundaries, data consistency, and operational observability.'
		WHEN 'frontend-engineer' THEN 'Build and improve frontend features with accessible UI architecture, state management, performance optimization, and maintainable tests.'
		WHEN 'fullstack-engineer' THEN 'Deliver end-to-end features across frontend and backend with aligned contracts, consistent UX, and safe rollout practices.'
		WHEN 'data-engineer' THEN 'Build reliable data pipelines and storage models with quality controls, schema evolution, and recoverability.'
		WHEN 'ml-ai-engineer' THEN 'Design and improve ML/AI systems with robust data/model lifecycle practices, evaluation criteria, and safe deployment behavior.'
		WHEN 'mobile-engineer' THEN 'Build and maintain mobile features with strong UX, platform-aware behavior, performance efficiency, and offline resilience.'
		WHEN 'devops-engineer' THEN 'Improve delivery pipelines, deployment workflows, and infrastructure automation with safe and repeatable operations.'
		WHEN 'sre-engineer' THEN 'Improve service reliability through incident response, SLO-driven prioritization, and operational hardening.'
		WHEN 'platform-engineer' THEN 'Improve internal platform capabilities and developer experience through reusable standards, tooling, and golden paths.'
		WHEN 'software-architect' THEN 'Guide system design decisions with clear boundaries, trade-off analysis, and migration strategies for scalable software evolution.'
		WHEN 'qa-automation-engineer' THEN 'Design and maintain reliable automated tests across unit, integration, and end-to-end layers with stable execution and actionable failures.'
		WHEN 'security-engineer' THEN 'Assess and improve application security through threat modeling, secure defaults, and vulnerability mitigation.'
		WHEN 'creative-director' THEN 'Lead high-level creative direction across storytelling, visuals, motion, and production choices.'
		WHEN 'designer' THEN 'Design visual and interaction solutions with user-centered clarity and coherent brand language.'
		WHEN 'manager' THEN 'Plan, prioritize, and coordinate teams and projects to deliver outcomes under constraints.'
		WHEN 'administrator' THEN 'Organize administrative operations, documentation, and process standardization efficiently.'
		WHEN 'mathematician' THEN 'Solve mathematical problems with clear derivations, proofs, and practical interpretations.'
		WHEN 'physicist' THEN 'Explain and apply physics concepts with equations, assumptions, and dimensional consistency.'
		WHEN 'biologist' THEN 'Analyze biological systems from molecular to ecosystem levels using evidence-based reasoning.'
		WHEN 'chemist' THEN 'Reason about chemical systems, reactions, and laboratory context with scientific rigor.'
		WHEN 'physician' THEN 'Provide medically oriented educational guidance with safety-first framing and clinical reasoning structure.'
		WHEN 'aeronautics-engineer' THEN 'Solve aeronautics problems spanning aircraft performance, stability, and propulsion trade-offs.'
		WHEN 'astronomer' THEN 'Interpret astronomical phenomena with observational context and physical principles.'
		WHEN 'musician' THEN 'Support composition, harmony, rhythm, arrangement, and musical analysis.'
		WHEN 'singer' THEN 'Guide vocal practice, interpretation, and performance preparation with safe vocal techniques.'
		WHEN 'philosopher' THEN 'Structure arguments, examine assumptions, and compare schools of thought rigorously.'
		WHEN 'marketing-sales' THEN 'Design integrated marketing and sales strategies grounded in funnel metrics and positioning.'
		WHEN 'seo-sem-specialist' THEN 'Optimize organic and paid search performance using keyword intent, content structure, and campaign tuning.'
		WHEN 'productivity-workflows' THEN 'Coordinate GitHub, Notion, Obsidian vault, and transactional email with explicit least-privilege tools.'
	END
FROM assistant_personalities
WHERE reference != 'general';

-- Spanish translations
INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'es', 'General', 'Asistente equilibrado para tareas del dia a dia y soporte general.'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'es',
	CASE reference
		WHEN 'backend-engineer' THEN 'Ingeniero Backend'
		WHEN 'frontend-engineer' THEN 'Ingeniero Frontend'
		WHEN 'fullstack-engineer' THEN 'Ingeniero Fullstack'
		WHEN 'data-engineer' THEN 'Ingeniero de Datos'
		WHEN 'ml-ai-engineer' THEN 'Ingeniero ML/IA'
		WHEN 'mobile-engineer' THEN 'Ingeniero Mobile'
		WHEN 'devops-engineer' THEN 'Ingeniero DevOps'
		WHEN 'sre-engineer' THEN 'Ingeniero SRE'
		WHEN 'platform-engineer' THEN 'Ingeniero de Plataforma'
		WHEN 'software-architect' THEN 'Arquitecto de Software'
		WHEN 'qa-automation-engineer' THEN 'Ingeniero QA Automatizado'
		WHEN 'security-engineer' THEN 'Ingeniero de Seguridad'
		WHEN 'creative-director' THEN 'Director Creativo'
		WHEN 'designer' THEN 'Diseñador'
		WHEN 'manager' THEN 'Gestor'
		WHEN 'administrator' THEN 'Administrador'
		WHEN 'mathematician' THEN 'Matemático'
		WHEN 'physicist' THEN 'Físico'
		WHEN 'biologist' THEN 'Biólogo'
		WHEN 'chemist' THEN 'Químico'
		WHEN 'physician' THEN 'Médico'
		WHEN 'aeronautics-engineer' THEN 'Ingeniero Aeronautico'
		WHEN 'astronomer' THEN 'Astrónomo'
		WHEN 'musician' THEN 'Músico'
		WHEN 'singer' THEN 'Cantante'
		WHEN 'philosopher' THEN 'Filósofo'
		WHEN 'marketing-sales' THEN 'Marketing y Ventas'
		WHEN 'seo-sem-specialist' THEN 'Especialista SEO/SEM'
		WHEN 'productivity-workflows' THEN 'Flujos de Productividad'
	END,
	NULL
FROM assistant_personalities
WHERE reference != 'general';

-- French translations
INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'fr', 'Général', 'Assistant équilibré pour les tâches quotidiennes et le support général.'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'fr',
	CASE reference
		WHEN 'backend-engineer' THEN 'Ingénieur Backend'
		WHEN 'frontend-engineer' THEN 'Ingénieur Frontend'
		WHEN 'fullstack-engineer' THEN 'Ingénieur Fullstack'
		WHEN 'data-engineer' THEN 'Ingénieur Données'
		WHEN 'ml-ai-engineer' THEN 'Ingénieur ML/IA'
		WHEN 'mobile-engineer' THEN 'Ingénieur Mobile'
		WHEN 'devops-engineer' THEN 'Ingénieur DevOps'
		WHEN 'sre-engineer' THEN 'Ingénieur SRE'
		WHEN 'platform-engineer' THEN 'Ingénieur Plateforme'
		WHEN 'software-architect' THEN 'Architecte Logiciel'
		WHEN 'qa-automation-engineer' THEN 'Ingénieur QA Automatisation'
		WHEN 'security-engineer' THEN 'Ingénieur Sécurité'
		WHEN 'creative-director' THEN 'Directeur Créatif'
		WHEN 'designer' THEN 'Concepteur'
		WHEN 'manager' THEN 'Manager'
		WHEN 'administrator' THEN 'Administrateur'
		WHEN 'mathematician' THEN 'Mathématicien'
		WHEN 'physicist' THEN 'Physicien'
		WHEN 'biologist' THEN 'Biologiste'
		WHEN 'chemist' THEN 'Chimiste'
		WHEN 'physician' THEN 'Médecin'
		WHEN 'aeronautics-engineer' THEN 'Ingénieur Aéronautique'
		WHEN 'astronomer' THEN 'Astronome'
		WHEN 'musician' THEN 'Musicien'
		WHEN 'singer' THEN 'Chanteur'
		WHEN 'philosopher' THEN 'Philosophe'
		WHEN 'marketing-sales' THEN 'Marketing et Ventes'
		WHEN 'seo-sem-specialist' THEN 'Spécialiste SEO/SEM'
		WHEN 'productivity-workflows' THEN 'Flux de Productivité'
	END,
	NULL
FROM assistant_personalities
WHERE reference != 'general';

-- German translations
INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'de', 'Allgemein', 'Ausgewogener Assistent für tägliche Aufgaben und allgemeine Unterstützung.'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'de',
	CASE reference
		WHEN 'backend-engineer' THEN 'Backend-Ingenieur'
		WHEN 'frontend-engineer' THEN 'Frontend-Ingenieur'
		WHEN 'fullstack-engineer' THEN 'Fullstack-Ingenieur'
		WHEN 'data-engineer' THEN 'Datentechnik-Ingenieur'
		WHEN 'ml-ai-engineer' THEN 'ML/KI-Ingenieur'
		WHEN 'mobile-engineer' THEN 'Mobile-Ingenieur'
		WHEN 'devops-engineer' THEN 'DevOps-Ingenieur'
		WHEN 'sre-engineer' THEN 'SRE-Ingenieur'
		WHEN 'platform-engineer' THEN 'Plattform-Ingenieur'
		WHEN 'software-architect' THEN 'Softwarearchitekt'
		WHEN 'qa-automation-engineer' THEN 'QA-Automatisierungs-Ingenieur'
		WHEN 'security-engineer' THEN 'Sicherheits-Ingenieur'
		WHEN 'creative-director' THEN 'Kreativdirektor'
		WHEN 'designer' THEN 'Designer'
		WHEN 'manager' THEN 'Manager'
		WHEN 'administrator' THEN 'Administrator'
		WHEN 'mathematician' THEN 'Mathematiker'
		WHEN 'physicist' THEN 'Physiker'
		WHEN 'biologist' THEN 'Biologe'
		WHEN 'chemist' THEN 'Chemiker'
		WHEN 'physician' THEN 'Arzt'
		WHEN 'aeronautics-engineer' THEN 'Luftfahrtingenieur'
		WHEN 'astronomer' THEN 'Astronom'
		WHEN 'musician' THEN 'Musiker'
		WHEN 'singer' THEN 'Sänger'
		WHEN 'philosopher' THEN 'Philosoph'
		WHEN 'marketing-sales' THEN 'Marketing und Vertrieb'
		WHEN 'seo-sem-specialist' THEN 'SEO/SEM-Spezialist'
		WHEN 'productivity-workflows' THEN 'Produktivitäts-Workflows'
	END,
	NULL
FROM assistant_personalities
WHERE reference != 'general';

-- Portuguese translations
INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'pt', 'Geral', 'Assistente equilibrado para tarefas diárias e suporte geral.'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'pt',
	CASE reference
		WHEN 'backend-engineer' THEN 'Engenheiro Backend'
		WHEN 'frontend-engineer' THEN 'Engenheiro Frontend'
		WHEN 'fullstack-engineer' THEN 'Engenheiro Fullstack'
		WHEN 'data-engineer' THEN 'Engenheiro de Dados'
		WHEN 'ml-ai-engineer' THEN 'Engenheiro ML/IA'
		WHEN 'mobile-engineer' THEN 'Engenheiro Mobile'
		WHEN 'devops-engineer' THEN 'Engenheiro DevOps'
		WHEN 'sre-engineer' THEN 'Engenheiro SRE'
		WHEN 'platform-engineer' THEN 'Engenheiro de Plataforma'
		WHEN 'software-architect' THEN 'Arquiteto de Software'
		WHEN 'qa-automation-engineer' THEN 'Engenheiro QA Automação'
		WHEN 'security-engineer' THEN 'Engenheiro de Segurança'
		WHEN 'creative-director' THEN 'Diretor Criativo'
		WHEN 'designer' THEN 'Designer'
		WHEN 'manager' THEN 'Gerenciador'
		WHEN 'administrator' THEN 'Administrador'
		WHEN 'mathematician' THEN 'Matemático'
		WHEN 'physicist' THEN 'Físico'
		WHEN 'biologist' THEN 'Biólogo'
		WHEN 'chemist' THEN 'Químico'
		WHEN 'physician' THEN 'Médico'
		WHEN 'aeronautics-engineer' THEN 'Engenheiro Aeronáutico'
		WHEN 'astronomer' THEN 'Astrônomo'
		WHEN 'musician' THEN 'Músico'
		WHEN 'singer' THEN 'Cantor'
		WHEN 'philosopher' THEN 'Filósofo'
		WHEN 'marketing-sales' THEN 'Marketing e Vendas'
		WHEN 'seo-sem-specialist' THEN 'Especialista SEO/SEM'
		WHEN 'productivity-workflows' THEN 'Fluxos de Produtividade'
	END,
	NULL
FROM assistant_personalities
WHERE reference != 'general';

-- Italian translations
INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'it', 'Generale', 'Assistente equilibrato per compiti quotidiani e supporto generale.'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'it',
	CASE reference
		WHEN 'backend-engineer' THEN 'Ingegnere Backend'
		WHEN 'frontend-engineer' THEN 'Ingegnere Frontend'
		WHEN 'fullstack-engineer' THEN 'Ingegnere Fullstack'
		WHEN 'data-engineer' THEN 'Ingegnere Dati'
		WHEN 'ml-ai-engineer' THEN 'Ingegnere ML/IA'
		WHEN 'mobile-engineer' THEN 'Ingegnere Mobile'
		WHEN 'devops-engineer' THEN 'Ingegnere DevOps'
		WHEN 'sre-engineer' THEN 'Ingegnere SRE'
		WHEN 'platform-engineer' THEN 'Ingegnere Piattaforma'
		WHEN 'software-architect' THEN 'Architetto Software'
		WHEN 'qa-automation-engineer' THEN 'Ingegnere QA Automazione'
		WHEN 'security-engineer' THEN 'Ingegnere Sicurezza'
		WHEN 'creative-director' THEN 'Direttore Creativo'
		WHEN 'designer' THEN 'Designer'
		WHEN 'manager' THEN 'Manager'
		WHEN 'administrator' THEN 'Amministratore'
		WHEN 'mathematician' THEN 'Matematico'
		WHEN 'physicist' THEN 'Fisico'
		WHEN 'biologist' THEN 'Biologo'
		WHEN 'chemist' THEN 'Chimico'
		WHEN 'physician' THEN 'Medico'
		WHEN 'aeronautics-engineer' THEN 'Ingegnere Aeronautico'
		WHEN 'astronomer' THEN 'Astronomo'
		WHEN 'musician' THEN 'Musicista'
		WHEN 'singer' THEN 'Cantante'
		WHEN 'philosopher' THEN 'Filosofo'
		WHEN 'marketing-sales' THEN 'Marketing e Vendite'
		WHEN 'seo-sem-specialist' THEN 'Specialista SEO/SEM'
		WHEN 'productivity-workflows' THEN 'Flussi di Produttività'
	END,
	NULL
FROM assistant_personalities
WHERE reference != 'general';

-- Russian translations
INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'ru', 'Общий', 'Сбалансированный ассистент для повседневных задач и общей поддержки.'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'ru',
	CASE reference
		WHEN 'backend-engineer' THEN 'Backend-инженер'
		WHEN 'frontend-engineer' THEN 'Frontend-инженер'
		WHEN 'fullstack-engineer' THEN 'Fullstack-инженер'
		WHEN 'data-engineer' THEN 'Инженер данных'
		WHEN 'ml-ai-engineer' THEN 'Инженер ML/ИИ'
		WHEN 'mobile-engineer' THEN 'Mobile-инженер'
		WHEN 'devops-engineer' THEN 'DevOps-инженер'
		WHEN 'sre-engineer' THEN 'SRE-инженер'
		WHEN 'platform-engineer' THEN 'Инженер платформы'
		WHEN 'software-architect' THEN 'Архитектор ПО'
		WHEN 'qa-automation-engineer' THEN 'Инженер QA автоматизации'
		WHEN 'security-engineer' THEN 'Инженер безопасности'
		WHEN 'creative-director' THEN 'Творческий директор'
		WHEN 'designer' THEN 'Дизайнер'
		WHEN 'manager' THEN 'Менеджер'
		WHEN 'administrator' THEN 'Администратор'
		WHEN 'mathematician' THEN 'Математик'
		WHEN 'physicist' THEN 'Физик'
		WHEN 'biologist' THEN 'Биолог'
		WHEN 'chemist' THEN 'Химик'
		WHEN 'physician' THEN 'Врач'
		WHEN 'aeronautics-engineer' THEN 'Инженер аэронавтики'
		WHEN 'astronomer' THEN 'Астроном'
		WHEN 'musician' THEN 'Музыкант'
		WHEN 'singer' THEN 'Певец'
		WHEN 'philosopher' THEN 'Философ'
		WHEN 'marketing-sales' THEN 'Маркетинг и продажи'
		WHEN 'seo-sem-specialist' THEN 'Специалист SEO/SEM'
		WHEN 'productivity-workflows' THEN 'Рабочие процессы продуктивности'
	END,
	NULL
FROM assistant_personalities
WHERE reference != 'general';

-- Japanese translations
INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'ja', '汎用', '日常的なタスクと一般サポート用のバランスの取れたアシスタント。'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'ja',
	CASE reference
		WHEN 'backend-engineer' THEN 'バックエンドエンジニア'
		WHEN 'frontend-engineer' THEN 'フロントエンドエンジニア'
		WHEN 'fullstack-engineer' THEN 'フルスタックエンジニア'
		WHEN 'data-engineer' THEN 'データエンジニア'
		WHEN 'ml-ai-engineer' THEN 'ML/AIエンジニア'
		WHEN 'mobile-engineer' THEN 'モバイルエンジニア'
		WHEN 'devops-engineer' THEN 'DevOpsエンジニア'
		WHEN 'sre-engineer' THEN 'SREエンジニア'
		WHEN 'platform-engineer' THEN 'プラットフォームエンジニア'
		WHEN 'software-architect' THEN 'ソフトウェアアーキテクト'
		WHEN 'qa-automation-engineer' THEN 'QA自動化エンジニア'
		WHEN 'security-engineer' THEN 'セキュリティエンジニア'
		WHEN 'creative-director' THEN 'クリエイティブディレクター'
		WHEN 'designer' THEN 'デザイナー'
		WHEN 'manager' THEN 'マネージャー'
		WHEN 'administrator' THEN '管理者'
		WHEN 'mathematician' THEN '数学者'
		WHEN 'physicist' THEN '物理学者'
		WHEN 'biologist' THEN '生物学者'
		WHEN 'chemist' THEN '化学者'
		WHEN 'physician' THEN '医師'
		WHEN 'aeronautics-engineer' THEN '航空工学エンジニア'
		WHEN 'astronomer' THEN '天文学者'
		WHEN 'musician' THEN 'ミュージシャン'
		WHEN 'singer' THEN '歌手'
		WHEN 'philosopher' THEN '哲学者'
		WHEN 'marketing-sales' THEN 'マーケティング&営業'
		WHEN 'seo-sem-specialist' THEN 'SEO/SEM専門家'
		WHEN 'productivity-workflows' THEN '生産性ワークフロー'
	END,
	NULL
FROM assistant_personalities
WHERE reference != 'general';

-- Chinese translations
INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'zh', '通用', '用于日常任务和一般支持的均衡助手。'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'zh',
	CASE reference
		WHEN 'backend-engineer' THEN '后端工程师'
		WHEN 'frontend-engineer' THEN '前端工程师'
		WHEN 'fullstack-engineer' THEN '全栈工程师'
		WHEN 'data-engineer' THEN '数据工程师'
		WHEN 'ml-ai-engineer' THEN 'ML/AI工程师'
		WHEN 'mobile-engineer' THEN '移动工程师'
		WHEN 'devops-engineer' THEN 'DevOps工程师'
		WHEN 'sre-engineer' THEN 'SRE工程师'
		WHEN 'platform-engineer' THEN '平台工程师'
		WHEN 'software-architect' THEN '软件架构师'
		WHEN 'qa-automation-engineer' THEN 'QA自动化工程师'
		WHEN 'security-engineer' THEN '安全工程师'
		WHEN 'creative-director' THEN '创意总监'
		WHEN 'designer' THEN '设计师'
		WHEN 'manager' THEN '经理'
		WHEN 'administrator' THEN '管理员'
		WHEN 'mathematician' THEN '数学家'
		WHEN 'physicist' THEN '物理学家'
		WHEN 'biologist' THEN '生物学家'
		WHEN 'chemist' THEN '化学家'
		WHEN 'physician' THEN '医生'
		WHEN 'aeronautics-engineer' THEN '航空工程师'
		WHEN 'astronomer' THEN '天文学家'
		WHEN 'musician' THEN '音乐家'
		WHEN 'singer' THEN '歌手'
		WHEN 'philosopher' THEN '哲学家'
		WHEN 'marketing-sales' THEN '营销和销售'
		WHEN 'seo-sem-specialist' THEN 'SEO/SEM专家'
		WHEN 'productivity-workflows' THEN '生产力工作流'
	END,
	NULL
FROM assistant_personalities
WHERE reference != 'general';

-- Hindi translations
INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'hi', 'सामान्य', 'दैनिक कार्यों और सामान्य समर्थन के लिए संतुलित सहायक।'
FROM assistant_personalities
WHERE reference = 'general';

INSERT OR IGNORE INTO assistant_personality_translations (personality_id, locale, display_name, description)
SELECT id, 'hi',
	CASE reference
		WHEN 'backend-engineer' THEN 'बैकएंड इंजीनियर'
		WHEN 'frontend-engineer' THEN 'फ्रंटएंड इंजीनियर'
		WHEN 'fullstack-engineer' THEN 'फुलस्टैक इंजीनियर'
		WHEN 'data-engineer' THEN 'डेटा इंजीनियर'
		WHEN 'ml-ai-engineer' THEN 'ML/AI इंजीनियर'
		WHEN 'mobile-engineer' THEN 'मोबाइल इंजीनियर'
		WHEN 'devops-engineer' THEN 'डेवऑप्स इंजीनियर'
		WHEN 'sre-engineer' THEN 'SRE इंजीनियर'
		WHEN 'platform-engineer' THEN 'प्लेटफॉर्म इंजीनियर'
		WHEN 'software-architect' THEN 'सॉफ्टवेयर आर्किटेक्ट'
		WHEN 'qa-automation-engineer' THEN 'QA ऑटोमेशन इंजीनियर'
		WHEN 'security-engineer' THEN 'सुरक्षा इंजीनियर'
		WHEN 'creative-director' THEN 'क्रिएटिव निर्देशक'
		WHEN 'designer' THEN 'डिजाइनर'
		WHEN 'manager' THEN 'प्रबंधक'
		WHEN 'administrator' THEN 'प्रशासक'
		WHEN 'mathematician' THEN 'गणितज्ञ'
		WHEN 'physicist' THEN 'भौतिकविद्'
		WHEN 'biologist' THEN 'जीवविज्ञानी'
		WHEN 'chemist' THEN 'रसायनज्ञ'
		WHEN 'physician' THEN 'चिकित्सक'
		WHEN 'aeronautics-engineer' THEN 'एयरोनॉटिक्स इंजीनियर'
		WHEN 'astronomer' THEN 'खगोलविद्'
		WHEN 'musician' THEN 'संगीतकार'
		WHEN 'singer' THEN 'गायक'
		WHEN 'philosopher' THEN 'दार्शनिक'
		WHEN 'marketing-sales' THEN 'मार्केटिंग और बिक्रय'
		WHEN 'seo-sem-specialist' THEN 'SEO/SEM विशेषज्ञ'
		WHEN 'productivity-workflows' THEN 'उत्पादकता वर्कफ़्लो'
	END,
	NULL
FROM assistant_personalities
WHERE reference != 'general';

INSERT OR IGNORE INTO assistant_tones (reference, is_default, is_active)
VALUES
	('formal', 1, 1),
	('casual', 0, 1),
	('empatico', 0, 1),
	('directo', 0, 1),
	('motivacional', 0, 1),
	('humoristico', 0, 1);

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'en',
	CASE reference
		WHEN 'formal' THEN 'Formal and Professional'
		WHEN 'casual' THEN 'Friendly and Casual'
		WHEN 'empatico' THEN 'Empathetic and Close'
		WHEN 'directo' THEN 'Direct and Concise'
		WHEN 'motivacional' THEN 'Motivational and Energetic'
		WHEN 'humoristico' THEN 'Humorous and Relaxed'
	END,
	NULL
FROM assistant_tones;

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'es',
	CASE reference
		WHEN 'formal' THEN 'Formal y Profesional'
		WHEN 'casual' THEN 'Amigable y Casual'
		WHEN 'empatico' THEN 'Empatico y Cercano'
		WHEN 'directo' THEN 'Directo y Conciso'
		WHEN 'motivacional' THEN 'Motivacional y Energetico'
		WHEN 'humoristico' THEN 'Humoristico y Desenfadado'
	END,
	NULL
FROM assistant_tones;

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'fr',
	CASE reference
		WHEN 'formal' THEN 'Formel & Professionnel'
		WHEN 'casual' THEN 'Amical & Décontracté'
		WHEN 'empatico' THEN 'Empathique & Chaleureux'
		WHEN 'directo' THEN 'Direct & Concis'
		WHEN 'motivacional' THEN 'Motivant & Énergique'
		WHEN 'humoristico' THEN 'Humoristique & Détendu'
	END,
	NULL
FROM assistant_tones;

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'de',
	CASE reference
		WHEN 'formal' THEN 'Formell & Professionell'
		WHEN 'casual' THEN 'Freundlich & Locker'
		WHEN 'empatico' THEN 'Einfühlsam & Nah'
		WHEN 'directo' THEN 'Direkt & Prägnant'
		WHEN 'motivacional' THEN 'Motivierend & Energetisch'
		WHEN 'humoristico' THEN 'Humorvoll & Unbekümmert'
	END,
	NULL
FROM assistant_tones;

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'pt',
	CASE reference
		WHEN 'formal' THEN 'Formal e Profissional'
		WHEN 'casual' THEN 'Amigável e Casual'
		WHEN 'empatico' THEN 'Empático e Próximo'
		WHEN 'directo' THEN 'Direto e Conciso'
		WHEN 'motivacional' THEN 'Motivacional e Energético'
		WHEN 'humoristico' THEN 'Humorístico e Descontraído'
	END,
	NULL
FROM assistant_tones;

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'it',
	CASE reference
		WHEN 'formal' THEN 'Formale & Professionale'
		WHEN 'casual' THEN 'Amichevole & Informale'
		WHEN 'empatico' THEN 'Empatico & Vicino'
		WHEN 'directo' THEN 'Diretto & Conciso'
		WHEN 'motivacional' THEN 'Motivazionale & Energico'
		WHEN 'humoristico' THEN 'Umoristico & Spensierato'
	END,
	NULL
FROM assistant_tones;

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'ru',
	CASE reference
		WHEN 'formal' THEN 'Формальный и Профессиональный'
		WHEN 'casual' THEN 'Дружелюбный и Непринуждённый'
		WHEN 'empatico' THEN 'Эмпатичный и Близкий'
		WHEN 'directo' THEN 'Прямой и Краткий'
		WHEN 'motivacional' THEN 'Мотивационный и Энергичный'
		WHEN 'humoristico' THEN 'Юмористический и Лёгкий'
	END,
	NULL
FROM assistant_tones;

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'ja',
	CASE reference
		WHEN 'formal' THEN 'フォーマル&プロフェッショナル'
		WHEN 'casual' THEN 'フレンドリー&カジュアル'
		WHEN 'empatico' THEN '共感的&親しみやすい'
		WHEN 'directo' THEN '直接的&簡潔'
		WHEN 'motivacional' THEN 'モチベーショナル&エネルギッシュ'
		WHEN 'humoristico' THEN 'ユーモラス&軽快'
	END,
	NULL
FROM assistant_tones;

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'zh',
	CASE reference
		WHEN 'formal' THEN '正式与专业'
		WHEN 'casual' THEN '友好与随性'
		WHEN 'empatico' THEN '共情与亲近'
		WHEN 'directo' THEN '直接与简洁'
		WHEN 'motivacional' THEN '激励与充沛'
		WHEN 'humoristico' THEN '幽默与轻松'
	END,
	NULL
FROM assistant_tones;

INSERT OR IGNORE INTO assistant_tone_translations (tone_id, locale, display_name, description)
SELECT id, 'hi',
	CASE reference
		WHEN 'formal' THEN 'औपचारिक और पेशेवर'
		WHEN 'casual' THEN 'मैत्रीपूर्ण और सहज'
		WHEN 'empatico' THEN 'सहानुभूतिपूर्ण और आत्मीय'
		WHEN 'directo' THEN 'प्रत्यक्ष और संक्षिप्त'
		WHEN 'motivacional' THEN 'प्रेरणादायक और ऊर्जावान'
		WHEN 'humoristico' THEN 'हास्यपूर्ण और हल्का-फुल्का'
	END,
	NULL
FROM assistant_tones;
