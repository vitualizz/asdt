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
    sub: 'Especialistas de IA que se pasan el trabajo entre sí, un resumen a la vez. Tú diriges la ruta; el equipo recuerda lo que importa.',
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
      { id: 'ux-ui', name: 'Diseño UX/UI', desc: 'Diseña la experiencia: flujos, pantallas, componentes y accesibilidad.', command: '/asdt-ux-ui' },
    ],
    orchestrator: {
      id: 'orchestrator',
      name: 'Asesor de ruta',
      desc: 'Analiza lo que pedís y recomienda qué especialistas involucrar y en qué orden. Si le decís que sí, los corre de a uno y hace una pausa después de cada uno — o ejecutás vos cada comando.',
      command: '/asdt',
    },
  },
  terminal: {
    tabs: ['Nueva feature', 'PM, paso a paso', 'Seguridad · STRIDE', 'Tu turno_'],
    tryLabel: 'Pruébalo',
  },
  pipeline: {
    title: 'Cómo funciona',
    sub: 'Describí lo que necesitás — /asdt te dice qué especialistas usar y en qué orden. Cada uno lee los resúmenes de los anteriores.',
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
    a11yDesc: 'Diagrama que muestra siete especialistas en secuencia: Investigador, PM, Arquitecto, Developer, QA, Seguridad, UX/UI — cada uno le pasa su resumen al siguiente.',
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
    kbNote: 'cada especialista devuelve un resumen · el siguiente lo lee',
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
        asdt: 'Cada especialista construye sobre los resúmenes de los anteriores, y las decisiones clave se recuerdan entre sesiones. ASDT recuerda lo que importa.',
      },
      {
        chat: 'Cada prompt improvisa el proceso desde cero: a veces brillante, a veces a lo loco.',
        asdt: 'Cada especialista trabaja con pasos definidos — explorar, especificar, ejecutar. Siempre.',
      },
      {
        chat: 'Llegás rápido a un demo, pero sin historias, decisiones ni pruebas que lo respalden.',
        asdt: 'Quedan resultados reales — el código en git, las decisiones clave en la memoria — listos para retomar mañana.',
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
    artifactReads: 'Lee estos resúmenes',
    artifactWrites: 'Devuelve su resumen',
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
    memoryAndEngram: 'Scratchpad y memoria',
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
    chainsLabel: 'La petición es',
    chains: {
      plan: 'un plan',
      resume: 'retomar un plan',
      build: 'construirlo',
    },
    specialistSteps: {
      'researcher:knowledge-recall': {
        purpose: 'Recuerda descubrimientos previos y decisiones relacionadas antes de explorar',
        produces: 'contexto (inline)',
      },
      'researcher:platform-analysis': {
        purpose: 'Carga convenciones de plataforma y la superficie de diseño desde config.yaml',
        produces: 'contexto (inline)',
      },
      'researcher:discovery': {
        purpose: 'Encuadra el problema, diverge en direcciones candidatas, juzga factibilidad con evidencia y converge en una recomendación',
        produces: 'resumen researcher',
      },
      'pm:knowledge-recall': {
        purpose: 'Recuerda requisitos y decisiones de alcance previas del proyecto',
        produces: 'contexto (inline)',
      },
      'pm:platform-analysis': {
        purpose: 'Carga convenciones de plataforma y la superficie de diseño desde config.yaml',
        produces: 'contexto (inline)',
      },
      'pm:backlog': {
        purpose: 'Convierte la petición en historias en orden de entrega, alcance explícito, NFRs medibles y criterios de aceptación',
        produces: 'resumen pm',
      },
      'ux-ui:knowledge-recall': {
        purpose: 'Recuerda decisiones de UX y patrones de componentes previos',
        produces: 'contexto (inline)',
      },
      'ux-ui:platform-analysis': {
        purpose: 'Carga el design system detectado del proyecto, sus convenciones y su huella, más la superficie de diseño desde config.yaml',
        produces: 'contexto (inline)',
      },
      'ux-ui:ux-spec': {
        purpose: 'Convierte el requisito en flujos para la superficie de diseño del proyecto — bifurcaciones, estados y copy — mapeados a componentes existentes, con la accesibilidad que cada uno debe cumplir',
        produces: 'ux-ui/ux-spec — lo lee ui-design',
      },
      'ux-ui:ui-design': {
        purpose: 'Diseña la UI que necesitan los flujos — layout, jerarquía, tipografía y estados de cada pantalla, cada componente faltante parte por parte, y la base visual cuando no hay design system',
        produces: 'resumen ux-ui',
      },
      'architect:knowledge-recall': {
        purpose: 'Recuerda decisiones arquitectónicas previas relevantes a este cambio',
        produces: 'contexto (inline)',
      },
      'architect:platform-analysis': {
        purpose: 'Carga convenciones de plataforma y huella de diseño desde config.yaml',
        produces: 'contexto (inline)',
      },
      'architect:design': {
        purpose: 'Decide el enfoque frente a sus alternativas y diseña el modelo de datos y la superficie de API que se derivan',
        produces: 'resumen architect',
      },
      'developer:knowledge-recall': {
        purpose: 'Recuerda decisiones previas relevantes a este cambio',
        produces: 'contexto (inline)',
      },
      'developer:platform-analysis': {
        purpose: 'Carga convenciones de plataforma y la superficie de diseño desde config.yaml',
        produces: 'contexto (inline)',
      },
      'developer:explore': {
        purpose: 'Mapea cada lugar al que debe llegar el cambio y lo que ya existe para reutilizar, y nombra las preguntas abiertas',
        produces: 'developer/explore — lo lee spec',
      },
      'developer:spec': {
        purpose: 'Define alcance, criterios de aceptación numerados, el enfoque más chico que resuelve todo y los ficheros que implement puede tocar — tests incluidos',
        produces: 'resumen developer (Stage: spec) — se muestra en la compuerta y se le pasa a implement',
      },
      'developer:approve': {
        purpose: 'Te muestra el alcance, los criterios de aceptación, el enfoque y los ficheros exactos — cuáles son tests — y espera aprobar, ajustar o parar antes de escribir nada',
        produces: 'tu visto bueno (inline)',
      },
      'developer:implement': {
        purpose: 'Escribe el código y los tests que declara el spec aprobado, dentro de sus límites de edición; nunca los ejecuta',
        produces: 'resumen developer (Stage: implemented)',
      },
      'developer:verify': {
        purpose: 'Ofrece los comandos de verificación, los corre solo si decís que sí y hace como mucho dos rondas de arreglo sobre los mismos ficheros',
        produces: 'el resultado de la verificación, sumado al resumen developer',
      },
      'security:knowledge-recall': {
        purpose: 'Recuerda hallazgos, modelos de amenaza y mitigaciones previas',
        produces: 'contexto (inline)',
      },
      'security:platform-analysis': {
        purpose: 'Carga convenciones de plataforma y huella arquitectónica',
        produces: 'contexto (inline)',
      },
      'security:assess': {
        purpose: 'Mapea la superficie de ataque, aplica STRIDE sobre ella y cruza las amenazas que sobreviven con las categorías OWASP aplicables',
        produces: 'security/assess — lo lee harden',
      },
      'security:harden': {
        purpose: 'Convierte la evaluación en hallazgos priorizados con mitigaciones concretas y una checklist accionable',
        produces: 'resumen security',
      },
      'qa:knowledge-recall': {
        purpose: 'Recuerda hallazgos de calidad y decisiones de test previas',
        produces: 'contexto (inline)',
      },
      'qa:platform-analysis': {
        purpose: 'Carga convenciones de plataforma y la superficie de diseño desde config.yaml',
        produces: 'contexto (inline)',
      },
      'qa:test-plan': {
        purpose: 'Encuentra los huecos y casos límite que los criterios de aceptación no cubren, los vuelve casos de test y da un veredicto go/no-go',
        produces: 'resumen qa',
      },
    },
    pipelineFlows: {
      'full-feature': {
        title: 'Obtener una sugerencia de pipeline',
        steps: {
          orchestrate: {
            description: 'Pedile al orquestador — ASDT analiza el pedido, recomienda qué especialistas involucrar y en qué orden, y si le decís que sí los corre de a uno, con una pausa después de cada uno.',
          },
          pm: {
            description: 'El PM define el alcance y escribe historias con criterios de aceptación, que devuelve como su resumen.',
          },
          architect: {
            description: 'Architect lee el resumen del PM, diseña el flujo de tokens y los contratos de API y devuelve el suyo.',
          },
          developer: {
            description: 'Developer lee los resúmenes del PM y del Architect, especifica el magic link handler, espera tu aprobación, lo implementa y devuelve su resumen.',
          },
          security: {
            description: 'Security revisa el mecanismo de autenticación, ejecuta análisis STRIDE y OWASP y devuelve su resumen — hallazgos con sus mitigaciones y un checklist de hardening.',
          },
        },
      },
      'mid-pipeline': {
        title: 'Continuar a mitad del pipeline',
        steps: {
          developer: {
            description: 'Developer lee los resúmenes anteriores desde el scratchpad de la sesión. En una sesión nueva reconstruye a partir de git, de las decisiones clave que recuerda la memoria y de vos. Sin pasar contexto manualmente.',
          },
          qa: {
            description: 'QA lee el resumen del Developer y planifica las pruebas: huecos en los criterios de aceptación, edge cases, casos Given/When/Then y un veredicto go/no-go.',
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
        kbNote: 'Researcher devuelve una dirección recomendada en su resumen — la toma el PM después (o el Architect, si se salta el PM).',
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
        note: 'Cuando el alcance y la arquitectura ya están definidos.',
      },
      'security-audit-existing-feature': {
        title: 'Auditoría de seguridad sobre una feature existente',
        note: 'Corré en cualquier momento — no requiere que haya corrido un especialista antes.',
      },
      'design-new-ui-component': {
        title: 'Diseñar un nuevo componente de UI',
        note: 'Cuando necesitás flujos, pantallas y componentes diseñados, y accesibilidad, antes de que el developer empiece a codear.',
      },
      'validate-test-coverage': {
        title: 'Validar la cobertura de tests antes de publicar',
        note: 'Corré después de Developer — QA lee el resumen del Developer.',
      },
      'pickup-developer-existing-adr': {
        title: 'Retomar en Developer después de un ADR existente',
        note: 'Cuando PM y Architect corrieron en una sesión anterior. Developer construye sobre las decisiones recordadas y el código.',
        kbNote: 'No hace falta volver a correr PM o Architect — contale al Developer qué quedó definido, o apoyate en las decisiones que recuerda la memoria.',
      },
      'add-security-review-inflight': {
        title: 'Agregar una revisión de seguridad a un pipeline en curso',
        note: 'Corré Security en cualquier momento sin reiniciar el pipeline.',
        kbNote: 'Security lee los resúmenes del Architect y del Developer cuando existen.',
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
        oneLiner: 'Detecta tu stack y escribe los hechos del proyecto en .asdt/config.yaml.',
      },
      asdt: {
        title: 'Sugerencia de ruteo del pipeline',
        oneLiner: 'Analiza el pedido, recomienda qué especialistas involucrar y en qué orden, y los corre si le decís que sí.',
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
        produces: 'resumen pm — el cambio en una frase, historias en orden de entrega, alcance dentro/fuera, criterios de aceptación, riesgos',
        doNotUseWhen: 'Ya tenés un backlog entry claro — volver a correr PM regenera las historias desde cero',
      },
      architect: {
        teaser: 'Decide cómo encajan las piezas antes de escribir código.',
        invokeWhen: 'La solución toca límites de servicios, modelos de datos o contratos de API, o tiene dos enfoques técnicos viables que vale la pena documentar',
        produces: 'resumen architect — la decisión con las alternativas que descartó, modelo de datos, superficie de API',
        doNotUseWhen: 'Necesitás código de implementación, planes de pruebas o specs de UX — Architect produce decisiones, no código',
      },
      developer: {
        teaser: 'Convierte un diseño ya definido en un plan que aprobás y después en código real.',
        invokeWhen: 'La forma de la solución está definida y necesitás un plan de implementación o código de producción escrito en el repo',
        produces: 'resumen developer — el plan aprobado (Stage: spec) y después los archivos escritos y el resultado de la verificación',
        doNotUseWhen: 'Todavía no fijaste el alcance ni la arquitectura — Developer va a implementar contra requisitos ambiguos',
      },
      qa: {
        teaser: 'Convierte los criterios de aceptación en un plan de pruebas sistemático con un veredicto de go/no-go.',
        invokeWhen: 'El código está listo para revisión, existen AC pero no fueron validados, o necesitás cobertura sistemática de edge cases y límites',
        produces: 'resumen qa — huecos en los criterios de aceptación, edge cases, test cases Given/When/Then, veredicto go/no-go',
        doNotUseWhen: 'Querés código de test ejecutable — QA produce especificaciones de test, no código que corra',
      },
      security: {
        teaser: 'Busca riesgos de auth, datos e integraciones antes de que se publiquen.',
        invokeWhen: 'La feature toca auth, sesiones, PII, integraciones externas, webhooks, o nuevos endpoints públicos de API',
        produces: 'resumen security — hallazgos con severidad high/medium/low, cada uno con su mitigación concreta, más un checklist de hardening',
        doNotUseWhen: 'Querés código de implementación o decisiones arquitectónicas — Security solo produce findings y checklists',
      },
      'ux-ui': {
        teaser: 'Diseña flujos, pantallas y componentes antes de que arranque la implementación.',
        invokeWhen: 'Una pantalla nueva o una UI a nivel de feature necesita diseño antes de implementar, o hay que mapear flujos de usuario',
        produces: 'resumen ux-ui — flujos con sus estados y textos, cada pantalla y cada componente faltante diseñados, accesibilidad y una dirección visual cuando no hay design system',
        doNotUseWhen: 'La pantalla ya está construida — una spec de UX entregada después de implementar llega demasiado tarde para darle forma',
      },
      researcher: {
        teaser: 'Explora un problema difuso antes de comprometerse con una dirección.',
        invokeWhen: 'El problema es difuso o abierto — necesitás descubrimiento y enmarcado antes de poder escribir requisitos',
        produces: 'resumen researcher — encuadre del problema, una dirección recomendada, cada descartada con su razón, evidencia de factibilidad',
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
      subject: { desc: 'Un tema corto en kebab-case. Una nota más nueva sobre el mismo tema reemplaza a la anterior.' },
      kind: { desc: 'Qué guarda la nota: decision, convention, lesson o finding.' },
      source: { desc: 'Qué especialista la escribió, sobre qué cambio y cuándo. Si dos notas comparten tema, gana la fecha más reciente.' },
      project: { desc: 'El alcance del proyecto: el nombre del directorio más un hash corto de su remoto origin de git (o de su ruta, si no tiene), para que las consultas no se mezclen entre proyectos y todos los clones compartan un alcance.' },
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
