import type { UIStrings } from '@i18n/types'

export const es: UIStrings = {
  nav: {
    home: 'ASDT',
    docs: 'Docs',
    github: 'GitHub',
    langPickerLabel: 'Idioma',
    specialists: 'Especialistas',
    howItWorks: 'Cómo funciona',
    githubStarLabel: 'Star en GitHub',
  },
  a11y: {
    skipToContent: 'Ir al contenido',
    themeToggleLabel: 'Cambiar tema',
    langPickerLabel: 'Seleccionar idioma',
  },
  hero: {
    eyebrow: 'Para Claude Code y OpenCode',
    headline: 'Todo un',
    headlineGrad: 'equipo de software',
    headlineSuffix: ', en tu terminal',
    sub: 'Especialistas de IA que se pasan el trabajo entre sí sobre una base de conocimientos compartida. Tú diriges la ruta; el equipo no olvida nada.',
    cta: 'Empezar',
    secondaryCta: 'Dale una estrella',
    installLabel: 'Instalación — un solo comando',
    installCmd: 'curl -fsSL https://raw.githubusercontent.com/vitualizz/asdt/main/install.sh | bash',
    copyLabel: 'Copiar',
    copiedLabel: '¡Copiado!',
    copyErrorLabel: 'Error al copiar — seleccioná manualmente',
  },
  specialists: {
    kicker: '01 · El equipo',
    title: 'Conoce a los especialistas',
    sub: 'Roles enfocados, cada uno con su oficio. Invocá a cualquiera directamente, o pedile a /asdt que te sugiera la ruta.',
    advisorStrip: 'El punto de partida',
    items: [
      { id: 'researcher', name: 'Investigador', desc: 'Explora problemas difusos antes de que existan requisitos: ideación divergente, análisis de viabilidad y una dirección recomendada.', command: '/asdt-researcher' },
      { id: 'pm', name: 'Product Manager', desc: 'Convierte ideas sueltas en historias de usuario con alcance claro y criterios definidos.', command: '/asdt-pm' },
      { id: 'architect', name: 'Arquitecto', desc: 'Toma las decisiones técnicas: diseño del sistema, contratos de API y registros de decisión.', command: '/asdt-architect' },
      { id: 'developer', name: 'Developer', desc: 'Convierte especificaciones y diseños en código de producción, con plan de implementación.', command: '/asdt-developer' },
      { id: 'qa', name: 'QA Engineer', desc: 'Construye la red de seguridad: planes de prueba, criterios de aceptación y reportes de calidad.', command: '/asdt-qa' },
      { id: 'security', name: 'Seguridad', desc: 'Encuentra los huecos que un atacante vería primero: modelos de amenaza y endurecimiento.', command: '/asdt-security' },
      { id: 'ux-ui', name: 'Diseño UX/UI', desc: 'Da forma a la experiencia: flujos, componentes y accesibilidad.', command: '/asdt-ux-ui' },
    ],
    orchestrator: {
      id: 'orchestrator',
      name: 'Asesor de ruta',
      desc: 'Analiza lo que pedís y recomienda qué especialistas involucrar y en qué orden. Vos confirmás el plan y ejecutás cada comando.',
      command: '/asdt',
    },
  },
  terminal: {
    tabs: ['Nueva feature', 'PM, paso a paso', 'Seguridad · STRIDE', 'Tu turno_'],
    tryLabel: 'Pruébalo',
  },
  pipeline: {
    title: 'Cómo funciona',
    sub: 'Describí lo que necesitás — /asdt te dice qué especialistas usar y en qué orden. Cada uno lee los artefactos del anterior desde la base de conocimiento compartida.',
    nodes: [
      { id: 'researcher', label: 'Investigador' },
      { id: 'pm', label: 'PM' },
      { id: 'architect', label: 'Arquitecto' },
      { id: 'developer', label: 'Developer' },
      { id: 'qa', label: 'QA' },
      { id: 'security', label: 'Seguridad' },
      { id: 'ux-ui', label: 'UX/UI' },
    ],
    a11yTitle: 'Pipeline de especialistas ASDT',
    a11yDesc: 'Diagrama que muestra siete especialistas en secuencia: Investigador, PM, Arquitecto, Developer, QA, Seguridad, UX/UI — todos conectados a través de una base de conocimiento compartida.',
  },
  recipes: {
    kicker: '02 · Cómo funciona',
    title: 'No es un pipeline. Es un equipo que vos componés.',
    sub: '/asdt sugiere una ruta según la tarea — vos la confirmás, la reordenás o vas directo a un especialista.',
    tabs: ['Nueva feature', 'Hotfix', 'Auditoría de seguridad', 'Pase de diseño'],
    notes: [
      '/asdt añade espacios de trabajo con permisos por rol',
      '/asdt corrige el crash al paginar resultados vacíos',
      '/asdt-security revisa el manejo de sesiones',
      '/asdt rediseña el onboarding para móvil',
    ],
    kbNote: 'cada paso guarda en la base de conocimiento · el siguiente especialista lee automáticamente',
  },
  vs: {
    kicker: '03 · Por qué ASDT',
    title: 'Un chat improvisa. Un equipo entrega.',
    sub: 'Los asistentes de código son buenísimos. Con método y memoria, se convierten en un equipo.',
    chatHead: 'Todo en un chat',
    asdtHead: 'Con ASDT',
    items: [
      {
        chat: 'El contexto se diluye a medida que la conversación crece — lo que decidiste ayer ya no existe.',
        asdt: 'Cada decisión queda en una base de conocimientos que todo el equipo consulta. ASDT recuerda.',
      },
      {
        chat: 'Cada prompt improvisa el proceso desde cero: a veces brillante, a veces a lo loco.',
        asdt: 'Cada especialista trabaja con pasos definidos — explorar, especificar, ejecutar. Siempre.',
      },
      {
        chat: 'Llegás rápido a un demo, pero sin historias, decisiones ni pruebas que lo respalden.',
        asdt: 'Quedan artefactos reales — historias, decisiones, pruebas — listos para retomar mañana.',
      },
    ],
  },
  ctaBand: {
    title: 'Contratá al equipo con un comando',
    sub: 'Funciona con Claude Code y OpenCode. Open source, licencia MIT.',
  },
  footer: {
    tagline: 'Impulsado por IA. Dirigido por ti.',
    githubLabel: 'GitHub',
    docsLabel: 'Docs',
    licenseLabel: 'Licencia MIT',
    credit: 'Hecho en abierto por Lee Palacios — vitualizz · © 2026 ASDT',
  },
  docs: {
    fallbackNotice: 'Esta página solo está disponible en inglés.',
    fallbackNoticeLink: 'Ver en inglés',
    gettingStarted: 'Comenzando',
    specialists: 'Especialistas',
    commands: 'Comandos',
    userFlows: 'Flujos de usuario',
    onThisPage: 'En esta página',
    stepsTitle: 'Pasos que ejecuta',
    searchPlaceholder: 'Buscar en docs...',
    searchNoResults: 'Sin resultados.',
    searchLabel: 'Buscar documentación',
    complexity: 'Complejidad',
    riskSurface: 'Superficie de riesgo',
    notCalled: 'No se invoca en este nivel — el Developer lo maneja directamente.',
    notEligible: 'No elegible — vuelve a simple.',
    notAutoInvoked: 'No se auto-invoca. Sigue siendo invocable bajo demanda.',
    produces: 'Produce',
    pipelinePosition: 'Posición en el pipeline',
    artifactReads: 'Lee de la base de conocimientos',
    artifactWrites: 'Escribe a la base de conocimientos',
    flowStep: 'Paso',
    flowOf: 'de',
    flowTitle: 'Cómo se ejecuta',
    troubleshooting: 'Solución de problemas',
    recipes: 'Recetas',
    tutorial: 'Tutorial',
    aiAssistants: 'Asistentes de IA',
    configuration: 'Configuración',
    howItWorksTitle: 'Cómo funciona',
    specialistModel: 'Modelo de especialistas',
    memoryAndEngram: 'Base de conocimiento y memoria',
    contributing: 'Contribuir',
    development: 'Desarrollo',
    groupGettingStarted: 'Empezar',
    groupConcepts: 'Conceptos',
    subgroupOverview: 'Resumen',
    subgroupDetail: 'Detalle',
    groupReference: 'Referencia',
    sidebarAriaLabel: 'Documentación',
    pageNavAriaLabel: 'Navegación de páginas',
    prev: 'Anterior',
    next: 'Siguiente',
  },
  commandCard: {
    slashLabel: 'Comando',
    cliLabel: 'CLI',
  },
  modelPreset: {
    recommendedLabel: 'Recomendado',
    defaultLabel: 'Por defecto',
    customizeCta: {
      title: 'Personalizar por paso',
      desc: 'Configurá el modelo de cada paso de especialista vos mismo.',
    },
  },
  specialistComparison: {
    invokeWhen: 'Invocar cuando…',
    doNotUseWhen: 'No usar cuando…',
    exampleCommand: 'Comando de ejemplo',
    fullDetailsLink: 'Ver detalle →',
  },
  recipesBrowser: {
    legend: 'Filtrar por objetivo',
    empty: 'No hay recetas para este filtro todavía — probá "Todas".',
    showingPrefix: 'Mostrando',
    recipeWord: 'receta',
    recipeWordPlural: 'recetas',
    forWord: 'para',
  },
  data: {
    specialistSteps: {
      'researcher:context-recall': {
        purpose: 'Busca en la memoria organizacional descubrimientos previos, decisiones relacionadas y restricciones conocidas antes de idear',
        produces: 'contexto (en línea)',
      },
      'researcher:divergent-ideation': {
        purpose: 'Enmarca el problema y genera direcciones candidatas divergentes — deliberadamente generativo, nunca selectivo',
        produces: 'researcher/ideation',
      },
      'researcher:feasibility-scan': {
        purpose: 'Evalúa cada idea con un veredicto de factibilidad verde/amarillo/rojo, evidencia de respaldo y estimación de esfuerzo',
        produces: 'researcher/feasibility',
      },
      'researcher:discovery-brief': {
        purpose: 'Converge en UNA dirección recomendada con su justificación; los candidatos descartados alimentan la lista de fuera de alcance del PM',
        produces: 'researcher/discovery-brief',
      },
      'researcher:decision-preservation': {
        purpose: 'Preserva la dirección elegida como conocimiento organizacional permanente mediante el campo summary del brief',
        produces: 'resumen (en línea)',
      },
      'pm:feature-intake': {
        purpose: 'Convierte el pedido en bruto en un enunciado de problema estructurado — extrae problema, objetivo, stakeholders y marca ambigüedades',
        produces: 'pm/feature-intake',
      },
      'pm:user-stories': {
        purpose: 'Escribe historias de usuario con prioridades MoSCoW y de 1 a 3 criterios de aceptación de alto nivel por historia',
        produces: 'pm/user-stories',
      },
      'pm:scope-analysis': {
        purpose: 'Define límites explícitos de dentro/fuera de alcance, puntos de integración y señales de riesgo de alcance',
        produces: 'pm/scope-analysis',
      },
      'pm:prioritization': {
        purpose: 'Ordena las historias por dependencia y riesgo; mueve las historias Won\'t a diferidas con razones explícitas',
        produces: 'pm/prioritization',
      },
      'pm:backlog-entry': {
        purpose: 'Consolida todos los artefactos del PM en el backlog entry final con resumen ejecutivo y lista ordenada de historias',
        produces: 'pm/backlog-entry',
      },
      'architect:load-constraints': {
        purpose: 'Lee el contexto de la plataforma y clasifica las restricciones como HARD (no negociables), SOFT (preferencias) u OPPORTUNITIES',
        produces: 'architect/constraints-analysis',
      },
      'architect:evaluate-approaches': {
        purpose: 'Compara 2 o 3 enfoques viables; elige uno con justificación explícita; documenta por qué se rechazaron las alternativas',
        produces: 'architect/approaches',
      },
      'architect:decision-record': {
        purpose: 'Escribe el ADR: contexto, decisión, alternativas y consecuencias — incluyendo las negativas. Solo positivas = incompleto',
        produces: 'architect/adr',
      },
      'architect:system-design': {
        purpose: 'Define el modelo de datos, la superficie de API (método/entradas/errores), los límites de servicios y la secuencia del happy path',
        produces: 'architect/system-design',
      },
      'architect:risk-analysis': {
        purpose: 'Identifica los 3 a 5 riesgos principales (rendimiento, seguridad, fiabilidad, acoplamiento, migración) con mitigaciones concretas',
        produces: 'architect/risks',
      },
      'architect:technical-handoff': {
        purpose: 'Consolida todo el trabajo de arquitectura en dos artefactos finales para Developer y QA',
        produces: 'architectural-decision + system-design',
      },
      'developer:explore': {
        purpose: 'Lee los archivos y módulos afectados; mapea los patrones de nombres y las restricciones antes de diseñar nada',
        produces: 'developer/dev-exploration',
      },
      'developer:spec': {
        purpose: 'Define el límite de alcance, responde las preguntas abiertas y escribe criterios de aceptación en formato Given/When/Then',
        produces: 'developer/dev-spec',
      },
      'developer:design': {
        purpose: 'Elige el enfoque técnico, define el modelo de datos y la forma de la API, y lista las restricciones clave de implementación',
        produces: 'developer/dev-design',
      },
      'developer:tasks': {
        purpose: 'Divide la implementación en tareas atómicas (<2h cada una) ordenadas por dependencia, con estimaciones S/M/L',
        produces: 'developer/dev-tasks',
      },
      'developer:implement': {
        purpose: 'Escribe el código de cada tarea — modo solo-plan (snippets) o modo escritura (archivos reales en los destinos declarados)',
        produces: 'developer/dev-implementation',
      },
      'qa:load-requirements': {
        purpose: 'Extrae y normaliza los criterios de aceptación de los artefactos previos al formato Given/When/Then',
        produces: 'qa/ac-list',
      },
      'qa:ac-validation': {
        purpose: 'Revisa cada AC en cuanto a atomicidad, medibilidad, independencia, completitud y no ambigüedad — reescribe los que fallan',
        produces: 'qa/ac-gaps',
      },
      'qa:edge-case-analysis': {
        purpose: 'Descubre edge cases mediante valores límite, particiones de equivalencia, transiciones de estado, acceso concurrente y límites de permisos',
        produces: 'qa/edge-cases',
      },
      'qa:test-strategy': {
        purpose: 'Define la pirámide de tests: qué cubre cada nivel, qué no, la estrategia de datos de prueba y la tolerancia a la inestabilidad',
        produces: 'qa/test-strategy',
      },
      'qa:test-case-generation': {
        purpose: 'Escribe specs de test estructuradas (Given/When/Then) para el happy path, los ACs validados y los edge cases críticos/altos',
        produces: 'qa/test-cases',
      },
      'qa:quality-report': {
        purpose: 'Verifica la cobertura de ACs, calcula el porcentaje y emite el veredicto READY / READY WITH CAVEATS / BLOCKED',
        produces: 'test-plan',
      },
      'qa:review': {
        purpose: 'Veredicto go/no-go de despliegue — decisión holística de preparación para release que integra todos los hallazgos de QA antes de preservar el conocimiento',
        produces: 'qa/qa-review',
      },
      'security:threat-modeling': {
        purpose: 'Aplica STRIDE: Spoofing, Tampering, Repudiation, Information Disclosure, DoS, Elevation of Privilege',
        produces: 'security/stride-threats',
      },
      'security:attack-surface': {
        purpose: 'Mapea los puntos de entrada, los límites de confianza y los flujos de datos — verifica validación/sanitización/codificación en cada paso',
        produces: 'security/attack-surface',
      },
      'security:owasp-analysis': {
        purpose: 'Revisa las 10 categorías del OWASP Top 10 (A01–A10) como APPLICABLE/NOT APPLICABLE y MITIGATED/AT RISK',
        produces: 'security/owasp-findings',
      },
      'security:hardening-checklist': {
        purpose: 'Deduplica los hallazgos, prioriza por severidad y agrupa por esfuerzo: quick wins / medio / significativo',
        produces: 'security-findings + hardening-checklist',
      },
      'ux-ui:feature-brief': {
        purpose: 'Identifica al actor principal, define el problema central (no la solución) y establece de 3 a 5 criterios de éxito observables',
        produces: 'ux-ui/feature-brief',
      },
      'ux-ui:information-architecture': {
        purpose: 'Organiza el contenido en secciones, prioriza la divulgación inmediata frente a la progresiva y define la ruta de navegación',
        produces: 'ux-ui/ia',
      },
      'ux-ui:user-flows': {
        purpose: 'Mapea el happy path, el flujo de error y de 2 a 3 flujos de edge cases como pasos numerados desde la perspectiva del actor',
        produces: 'ux-ui/flows',
      },
      'ux-ui:component-mapping': {
        purpose: 'Clasifica cada estado de UI como reuse / extend / new — quality gate: se requiere una relación de reutilización >2:1',
        produces: 'ux-ui/components',
      },
      'ux-ui:responsive-strategy': {
        purpose: 'Specs de breakpoints mobile-first para componentes nuevos/extendidos; confirma áreas táctiles mínimas de 44×44px',
        produces: 'ux-ui/responsive',
      },
      'ux-ui:ux-handoff': {
        purpose: 'Consolida todo el trabajo de UX en ux-brief (flows + IA) y component-spec (inventario + props + events)',
        produces: 'ux-brief + component-spec',
      },
    },
    pipelineFlows: {
      'full-feature': {
        title: 'Obtener una sugerencia de pipeline',
        steps: {
          orchestrate: {
            description: 'Pedile al orquestador — ASDT analiza el pedido y recomienda qué especialistas involucrar y en qué orden.',
          },
          pm: {
            description: 'PM define el alcance, escribe historias de usuario con criterios de aceptación y guarda pm/backlog-entry en la base de conocimientos.',
          },
          architect: {
            description: 'Architect lee el backlog entry, diseña el flujo de tokens y los contratos de API, guarda architectural-decision + system-design.',
          },
          developer: {
            description: 'Developer lee el ADR y el system design, implementa el magic link handler y guarda dev-implementation.',
          },
          security: {
            description: 'Security revisa el mecanismo de autenticación, ejecuta análisis STRIDE y OWASP, guarda security-findings + hardening-checklist.',
          },
        },
      },
      'mid-pipeline': {
        title: 'Continuar a mitad del pipeline',
        steps: {
          developer: {
            description: 'Developer lee los artefactos previos de la base de conocimientos automáticamente, incluso de sesiones anteriores. Sin pasar contexto manualmente.',
          },
          qa: {
            description: 'QA carga dev-implementation y ejecuta su flujo completo: validación de ACs, análisis de edge cases y generación de tests.',
          },
        },
      },
    },
    recipes: {
      'ship-new-feature': {
        title: 'Lanzar una feature nueva de cara al usuario',
        note: 'Empezá acá cuando el pedido de una feature está en lenguaje vago y necesita la secuencia completa PM → Architect → Developer.',
        kbNote: 'Revisá la secuencia de especialistas sugerida y después corré cada uno en orden.',
      },
      'new-rest-endpoint': {
        title: 'Construir un nuevo endpoint REST',
        note: 'Para cambios solo de backend donde el contrato de la API es la decisión de diseño principal.',
      },
      'new-screen-ux-first': {
        title: 'Pantalla nueva con diseño UX primero',
        note: 'Para features de cara al usuario donde el diseño de UI debería preceder a la implementación.',
      },
      'security-review-before-shipping': {
        title: 'Feature con revisión de seguridad antes de publicar',
        note: 'Para cualquier cosa que toque auth, pagos, PII (información de identificación personal) o integraciones externas.',
      },
      'explore-before-planning': {
        title: 'Explorar antes de planificar (problema difuso)',
        note: 'Cuando el problema no está claro y necesitás descubrimiento antes de definir requisitos.',
        kbNote: 'Researcher produce un discovery brief con una dirección recomendada — pasáselo al PM después.',
      },
      'lock-scope-user-stories': {
        title: 'Fijar el alcance y escribir historias de usuario',
        note: 'Cuando tenés una idea de feature clara y solo necesitás requisitos estructurados.',
      },
      'document-architecture-decision': {
        title: 'Documentar una decisión de arquitectura',
        note: 'Cuando el enfoque técnico necesita un ADR formal y un diseño de sistema.',
      },
      'write-code-from-settled-spec': {
        title: 'Escribir código de producción a partir de una spec ya definida',
        note: 'Cuando el alcance y la arquitectura ya están definidos en la base de conocimiento.',
      },
      'security-audit-existing-feature': {
        title: 'Auditoría de seguridad sobre una feature existente',
        note: 'Corré en cualquier momento — no requiere que haya corrido un especialista antes.',
      },
      'design-new-ui-component': {
        title: 'Diseñar un nuevo componente de UI',
        note: 'Cuando necesitás una spec de componente antes de que el developer empiece a codear.',
      },
      'validate-test-coverage': {
        title: 'Validar la cobertura de tests antes de publicar',
        note: 'Corré después de Developer — QA lee la implementación automáticamente.',
      },
      'pickup-developer-existing-adr': {
        title: 'Retomar en Developer después de un ADR existente',
        note: 'Cuando PM y Architect corrieron en una sesión anterior. Developer carga los artefactos previos automáticamente.',
        kbNote: 'No hace falta volver a correr PM o Architect — los artefactos están en la base de conocimiento.',
      },
      'add-security-review-inflight': {
        title: 'Agregar una revisión de seguridad a un pipeline en curso',
        note: 'Corré Security en cualquier momento sin reiniciar el pipeline.',
        kbNote: 'Security lee system-design y dev-implementation de la base de conocimiento.',
      },
      'qa-completed-feature-no-pipeline': {
        title: 'Hacer QA de una feature terminada sin pipeline completo',
        note: 'Cuando se construyó una feature sin ASDT — corré QA contra el código existente.',
      },
    },
    recipeCategories: {
      all: { label: 'Todas' },
      'from-scratch': { label: 'Empezar de cero' },
      'add-to-existing': { label: 'Añadir a código existente' },
      'review-harden': { label: 'Revisar y endurecer' },
      'understand-document': { label: 'Entender y documentar' },
    },
    commands: {
      'asdt-tui': {
        title: 'Interfaz interactiva de terminal',
        oneLiner: 'La única herramienta CLI — revisa tu setup e instala o actualiza las skills de los especialistas ASDT.',
      },
      'asdt-init': {
        title: 'Inicializar ASDT',
        oneLiner: 'Detecta tu stack y crea .asdt/config.yaml y .asdt/knowledge/platform.yaml.',
      },
      asdt: {
        title: 'Sugerencia de ruteo del pipeline',
        oneLiner: 'Analiza el pedido y recomienda qué especialistas involucrar y en qué orden.',
      },
      'asdt-researcher': {
        title: 'Researcher',
        oneLiner: 'Corre solo al especialista Researcher — descubrimiento para problemas difusos antes de que existan requisitos.',
      },
      'asdt-pm': {
        title: 'Product Manager',
        oneLiner: 'Corre solo al especialista Product Manager.',
      },
      'asdt-architect': {
        title: 'Architect',
        oneLiner: 'Corre solo al especialista Architect.',
      },
      'asdt-developer': {
        title: 'Developer',
        oneLiner: 'Corre solo al especialista Developer.',
      },
      'asdt-qa': {
        title: 'QA',
        oneLiner: 'Corre solo al especialista QA.',
      },
      'asdt-security': {
        title: 'Security',
        oneLiner: 'Corre solo al especialista Security.',
      },
      'asdt-ux-ui': {
        title: 'UX/UI',
        oneLiner: 'Corre solo al especialista UX/UI.',
      },
    },
    specialistComparison: {
      pm: {
        teaser: 'Fija el alcance y convierte pedidos vagos en historias de usuario estructuradas.',
        invokeWhen: 'El pedido es vago o de cara al usuario, el alcance no está definido, o todavía no existen historias de usuario',
        produces: 'pm/backlog-entry — resumen de la feature, historias de usuario ordenadas con AC, alcance dentro/fuera, riesgos',
        doNotUseWhen: 'Ya tenés un backlog entry claro — volver a correr PM regenera las historias desde cero',
      },
      architect: {
        teaser: 'Decide cómo encajan las piezas antes de escribir código.',
        invokeWhen: 'La solución toca límites de servicios, modelos de datos o contratos de API, o tiene dos enfoques técnicos viables que vale la pena documentar',
        produces: 'architectural-decision (ADR) + system-design — modelo de datos, superficie de API, límites de servicios',
        doNotUseWhen: 'Necesitás código de implementación, planes de pruebas o specs de UX — Architect produce decisiones, no código',
      },
      developer: {
        teaser: 'Convierte un diseño ya definido en un plan de implementación ordenado o código real.',
        invokeWhen: 'La forma de la solución está definida y necesitás un plan de implementación ordenado o código de producción escrito en el repo',
        produces: 'developer/dev-implementation — manifiesto de archivos ordenado y plan de código',
        doNotUseWhen: 'Todavía no fijaste el alcance ni la arquitectura — Developer va a implementar contra requisitos ambiguos',
      },
      qa: {
        teaser: 'Convierte los criterios de aceptación en un plan de pruebas sistemático con un veredicto de go/no-go.',
        invokeWhen: 'El código está listo para revisión, existen AC pero no fueron validados, o necesitás cobertura sistemática de edge cases y límites',
        produces: 'test-plan — % de cobertura de AC, gaps sin cubrir, lista completa de test cases Given/When/Then, veredicto de calidad',
        doNotUseWhen: 'Querés código de test ejecutable — QA produce especificaciones de test, no código que corra',
      },
      security: {
        teaser: 'Busca riesgos de auth, datos e integraciones antes de que se publiquen.',
        invokeWhen: 'La feature toca auth, sesiones, PII, integraciones externas, webhooks, o nuevos endpoints públicos de API',
        produces: 'security-findings (con severidad y CWE) + hardening-checklist — qué arreglar sí o sí vs qué se puede postergar',
        doNotUseWhen: 'Querés código de implementación o decisiones arquitectónicas — Security solo produce findings y checklists',
      },
      'ux-ui': {
        teaser: 'Mapea flujos y componentes antes de que arranque la implementación.',
        invokeWhen: 'Una pantalla nueva o una UI a nivel de feature necesita diseño antes de implementar, o hay que mapear flujos de usuario',
        produces: 'ux-brief (flujos, IA, criterios de éxito) + component-spec — inventario de componentes reusados/extendidos/nuevos',
        doNotUseWhen: 'La pantalla ya está construida — una spec de UX entregada después de implementar llega demasiado tarde para darle forma',
      },
      researcher: {
        teaser: 'Explora un problema difuso antes de comprometerse con una dirección.',
        invokeWhen: 'El problema es difuso o abierto — necesitás descubrimiento y enmarcado antes de poder escribir requisitos',
        produces: 'researcher/discovery-brief — enmarcado del problema, conjunto divergente de ideas, análisis de factibilidad, una dirección recomendada',
        doNotUseWhen: 'Ya tenés un problema bien definido — Researcher explora; no produce historias de usuario ni ADRs',
      },
    },
    tutorialStages: {
      install: { label: 'Instalar' },
      initialize: { label: 'Inicializar' },
      recommendation: { label: 'Recomendación' },
      pm: { label: 'PM' },
      architect: { label: 'Arquitecto' },
      developer: { label: 'Developer' },
    },
    artifactAnatomy: {
      title: { desc: 'El nombre legible del artefacto.' },
      topicKey: { desc: 'La clave con la que se recupera automáticamente, sin matching difuso.' },
      type: { desc: 'La categoría del artefacto — architecture, decision, bugfix, etc.' },
      project: { desc: 'El proyecto al que pertenece, para que los resultados no se mezclen entre proyectos.' },
    },
    modelPresets: {
      chameleon: {
        label: 'Chameleon',
        desc: 'Mantiene el modelo que tu asistente ya tiene definido (quita el campo model: para que cada asistente use su propio valor por defecto).',
      },
      sprinter: {
        label: 'Sprinter',
        desc: 'El más rápido y económico en todo.',
      },
      craftsman: {
        label: 'Craftsman',
        desc: 'El balance recomendado entre velocidad y capacidad (los valores por defecto que se envían, tal cual).',
      },
      strategist: {
        label: 'Strategist',
        desc: 'Más capacidad para análisis y decisiones.',
      },
      mastermind: {
        label: 'Mastermind',
        desc: 'Máxima capacidad donde más importa.',
      },
    },
    artifactSentinels: {
      'Problem (raw)': { label: 'Problema (en bruto)' },
      'Request (raw)': { label: 'Solicitud (en bruto)' },
    },
  },
}
