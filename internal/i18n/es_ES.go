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
		TitleSelectProvider:   "Seleccionar Proveedor de Memoria",
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

		HintNavigate:       "navegar",
		HintSelect:         "seleccionar",
		HintQuit:           "salir",
		HintToggle:         "alternar",
		HintAllNone:        "todos/ninguno",
		HintBack:           "volver",
		HintBackToMenu:     "volver al menú",
		HintCycleMode:      "cambiar modo",
		HintCycleModel:     "cambiar modelo",
		HintExpand:         "expandir/contraer",
		HintResetDefault:   "restablecer valores",
		HintContinue:       "continuar",
		HintChecking:       "comprobando",
		HintEnvironment:    "entorno",
		HintEngramRequired: "obligatorio — ver arriba",

		HintGroupNav:      "Nav",
		HintGroupActions:  "Acciones",
		HintGroupStatus:   "Estado",
		HintGroupRequired: "Obligatorio",

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
		LabelProvider:   "Proveedor: ",
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

		SectionYourEnvironment: "Tu Entorno",
		SectionMemoryProvider:  "Proveedor de Memoria",
		SectionAIEnhancements:  "Mejoras de IA",

		PrefEngramRequired: "⚠  Engram es obligatorio para mantener la memoria entre sesiones de IA.",
		PrefEngramInstall:  "Instala: https://github.com/Gentleman-Programming/engram",
		PrefEngramRestart:  "Después de instalarlo, reinicia asdt-tui.",
	},
	Dashboard: DashboardStrings{
		LabelInstalled: "Instalado",
		LabelPersona:   "Persona",
	},
	Personas: PersonaStrings{
		Sky:         "Aguda y minuciosa. Critica las malas prácticas con ingenio y luego te lo explica todo.",
		Toffy:       "Cálida y entusiasta. Respuestas sencillas, sarcasmo accidental.",
		Atreus:      "Audaz y temerario. Lanza primero la solución poco convencional — y funciona.",
		Babi:        "Tu fan número uno. Cálida, rápida y lo celebra todo a lo grande.",
		LeePalacios: "Amante de los gatos, programador, otaku. Tu compañero de pair programming para todo el camino.",
	},
}
