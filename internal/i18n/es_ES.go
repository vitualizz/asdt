package i18n

// SpanishSpain is the es-ES peninsular Spanish catalog (Spain register: tú/
// vosotros, peninsular imperatives "pulsa"/"elige", casual "vale" where natural;
// no voseo and no Rioplatense forms).
//
// HAND-AUTHORED: this catalog is flagged for human/UX review — do not auto-accept
// the peninsular register without a native reviewer.
var SpanishSpain = Catalog{
	Installer: InstallerStrings{
		TitleMainMenu:         "asdt-tui",
		TitlePreflightCheck:   "Comprobación Previa",
		TitleSelectAssistants: "Seleccionar Asistentes",
		TitleSelectProvider:   "Seleccionar Proveedores",
		TitleModelGate:        "Modelos de IA",
		TitleModelSetup:       "Elegir Modelos por Paso",
		TitleAgentSetup:       "Personalidad del Agente",
		TitleEmojiPref:        "Preferencia de Emojis",
		TitleAgentWriteMode:   "Configuraciones Existentes Detectadas",
		TitleReview:           "Revisar y Confirmar",
		TitleInstalling:       "Instalando...",
		TitleDone:             "Instalación Completa",
		TitleDashboard:        "Panel",
		TitleLanguageSelect:   "Idioma",

		StepWord:   "paso",
		StepOfWord: "de",

		MenuInstall:   "Instalar / Actualizar Skills",
		MenuDashboard: "Panel",
		MenuQuit:      "Salir",

		HintNavigate:     "navegar",
		HintSelect:       "seleccionar",
		HintQuit:         "salir",
		HintToggle:       "alternar",
		HintAllNone:      "todos/ninguno",
		HintBack:         "volver",
		HintBackToMenu:   "volver al menú",
		HintCycleMode:    "cambiar modo",
		HintCycleModel:   "cambiar modelo",
		HintExpand:       "expandir/contraer",
		HintResetDefault: "restablecer valores",
		HintContinue:     "continuar",
		HintChecking:     "comprobando",
		HintEnvironment:  "entorno",

		HintGroupNav:     "Nav",
		HintGroupActions: "Acciones",
		HintGroupStatus:  "Estado",

		BtnContinue: "[ Continuar → ]",
		BtnSkip:     "[ Omitir → ]",
		BtnInstall:  "[ Instalar → ]",

		BodyAgentSetupSubtitle: "Configura cómo se comportan tus asistentes de IA en todas las herramientas",
		BodyEmojiPrefSubtitle:  "¿Quieres que tus asistentes usen emojis en sus respuestas?",
		BodyModelSetupSubtitle: "Cada paso de especialista se ejecuta en su propio modelo. Proveedores detectados: %s",

		OptionModelsPreset0:        "Camaleón — mantén el modelo que tu asistente ya tiene definido",
		OptionModelsPreset1:        "Velocista — lo más rápido y económico en todo",
		OptionModelsPreset2:        "Artesano — el equilibrio recomendado entre velocidad y capacidad",
		OptionModelsPreset3:        "Estratega — más capacidad para el análisis y las decisiones",
		OptionModelsPreset4:        "Mente Maestra — máxima capacidad donde más importa",
		OptionModelsPreset5:        "Personalizar por paso — elige el modelo de cada paso de especialista",
		HintPresetAssignment0:      "cada paso hereda el modelo propio de tu asistente",
		HintPresetAssignment1:      "decisiones → equilibrado · análisis → rápido · tareas ligeras → rápido",
		HintPresetAssignment2:      "decisiones → máximo · análisis → equilibrado · tareas ligeras → rápido",
		HintPresetAssignment3:      "decisiones → máximo · análisis → equilibrado · tareas ligeras → equilibrado",
		HintPresetAssignment4:      "decisiones → nivel máximo · análisis → máximo · tareas ligeras → equilibrado",
		HintPresetAssignment5:      "elige a mano un modelo para cada paso de especialista",
		HintRunsAs:                 "→ se ejecuta como %s en Claude Code",
		LabelModels:                "Modelos:",
		ValueModelsPreset:          "preajuste %s",
		ValueModelsCustomized:      "%d personalizados",
		BodyAgentWriteMode:         "Elige cómo gestionar cada configuración existente:",
		BodyInstalling:             "Instalando asistentes y skills...",
		BodyLanguageSelectSubtitle: "¿En qué idioma quieres usar el instalador?",

		OptionEmojiYes: "Sí — usar emojis",
		OptionEmojiNo:  "No — solo texto plano",

		WarnExistingConfig:     "⚠  Se ha detectado una configuración global de agente existente — se sobrescribirá.",
		WarnExistingConfigNote: "   Continúa bajo tu responsabilidad.",
		InfoWriteGlobal:        "Esto escribirá en tu configuración global de asistentes de IA.",

		LabelAssistants: "Asistentes:",
		LabelMemory:     "Memoria:   ",
		LabelCodeIntel:  "Código:    ",
		LabelPersona:    "Persona:   ",
		LabelEmojis:     "Emojis:    ",
		LabelConfig:     "Config:",
		LabelSkipped:    "omitido",

		LabelModeKeep:      "Mantener",
		LabelModeOverwrite: "Sobrescribir",
		LabelModeAppend:    "Añadir",

		LabelAgentConfig:      "Config del Agente:",
		MsgAgentConfigWritten: "configuración de %s escrita",
		MsgAgentConfigSkipped: "– %s: omitido (se mantuvo la configuración existente)",
		MsgOpenCodeNote:       "Nota: OpenCode lo usa como anulación global para todos los proyectos.",
		MsgGetStarted:         "Abre un proyecto y ejecuta /asdt para empezar.",
		MsgStaleRemoved:       "se han limpiado %d archivo(s) obsoleto(s)",

		SectionYourEnvironment:  "Tu Entorno",
		SectionMemoryProvider:   "Proveedor de Memoria",
		SectionCodeIntelligence: "Inteligencia de Código",

		OptionNoMemory:         "Sin memoria",
		DescNoMemory:           "ASDT funciona sin un proveedor; se recomienda uno porque conserva el conocimiento clave del proyecto entre sesiones.",
		WarnProviderMissing:    "⚠  No se ha encontrado %s. ASDT se instala vinculado a él y funciona sin memoria hasta que lo instales y lo conectes.",
		PrefProviderInstall:    "Instala: %s",
		DetailProviderNotFound: "no encontrado — opcional, recomendado: %s",
		ValueMemoryNone:        "ninguna (recomendado: %s, para conservar el conocimiento clave del proyecto entre sesiones)",
		ValueMemoryMissing:     "%s (no se ha encontrado — instala: %s)",

		OptionNoCodeIntel:       "Sin inteligencia de código",
		DescNoCodeIntel:         "Los especialistas buscan en el código con grep y lectura; un proveedor responde las preguntas de símbolos y estructura desde un índice.",
		WarnCodeIntelMissing:    "⚠  No se ha encontrado %s. ASDT se instala vinculado a él; los especialistas usan grep y lectura hasta que lo instales e indexes un proyecto.",
		DetailCodeIntelNotFound: "no encontrado — opcional: %s",
		ValueCodeIntelNone:      "ninguna (grep y lectura)",
		ValueCodeIntelMissing:   "%s (no se ha encontrado — instala: %s)",
	},
	Dashboard: DashboardStrings{
		LabelInstalled:   "Instalado",
		LabelPersona:     "Persona",
		LabelMemory:      "Memoria",
		ValueNoMemory:    "ninguna",
		LabelCodeIntel:   "Inteligencia de código",
		ValueNoCodeIntel: "ninguna",
	},
	Personas: PersonaStrings{
		Sky:         "Aguda y minuciosa. Critica las malas prácticas con ingenio y luego te lo explica todo.",
		Toffy:       "Cálida y entusiasta. Respuestas sencillas, sarcasmo accidental.",
		Atreus:      "Audaz y temerario. Lanza primero la solución poco convencional — y funciona.",
		Babi:        "Tu fan número uno. Cálida, rápida y lo celebra todo a lo grande.",
		LeePalacios: "Amante de los gatos, programador, otaku. Tu compañero de pair programming para todo el camino.",
	},
}
