export interface SearchEntry {
  path: string
  page: string
  group: string
  label: string
  desc?: string
}

export const searchIndex: SearchEntry[] = [
  { path: "/performance", page: "nav.performance", group: "perf.powerPlan", label: "perf.ultimatePerf", desc: "perf.ultimatePerfDesc" },
  { path: "/performance", page: "nav.performance", group: "perf.powerPlan", label: "perf.cState", desc: "perf.cStateDesc" },
  { path: "/performance", page: "nav.performance", group: "perf.memory", label: "perf.superfetch", desc: "perf.superfetchDesc" },
  { path: "/performance", page: "nav.performance", group: "perf.memory", label: "perf.memCompress", desc: "perf.memCompressDesc" },
  { path: "/performance", page: "nav.performance", group: "perf.display", label: "perf.fullscreen", desc: "perf.fullscreenDesc" },
  { path: "/performance", page: "nav.performance", group: "perf.display", label: "perf.windowed", desc: "perf.windowedDesc" },
  { path: "/performance", page: "nav.performance", group: "perf.display", label: "perf.mpo", desc: "perf.mpoDesc" },
  { path: "/performance", page: "nav.performance", group: "perf.display", label: "perf.clearShader", desc: "perf.clearShaderDesc" },
  { path: "/performance", page: "nav.performance", group: "perf.xboxServices", label: "perf.gameBar", desc: "perf.gameBarDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.desktop", label: "pers.notification", desc: "pers.notificationDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.desktop", label: "pers.legacyBalloon", desc: "pers.legacyBalloonDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.desktop", label: "pers.edgeSwipe", desc: "pers.edgeSwipeDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.desktop", label: "pers.contextMenu", desc: "pers.contextMenuDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.desktop", label: "pers.win11TaskManager", desc: "pers.win11TaskManagerDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.desktop", label: "pers.widgets", desc: "pers.widgetsDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.desktop", label: "pers.explorerHome", desc: "pers.explorerHomeDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.desktop", label: "pers.explorerGallery", desc: "pers.explorerGalleryDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.appearance", label: "pers.removeArrow", desc: "pers.removeArrowDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.appearance", label: "pers.removeText", desc: "pers.removeTextDesc" },
  { path: "/personalization", page: "nav.personalization", group: "pers.appearance", label: "pers.removeShield", desc: "pers.removeShieldDesc" },
  { path: "/security", page: "nav.security", group: "sec.defenderModule", label: "sec.defenderTitle" },
  { path: "/security", page: "nav.security", group: "sec.defenderModule", label: "sec.defenderTpTitle" },
  { path: "/security", page: "nav.security", group: "sec.systemProtection", label: "sec.uac", desc: "sec.uacDesc" },
  { path: "/security", page: "nav.security", group: "sec.systemProtection", label: "sec.vbs", desc: "sec.vbsDesc" },
  { path: "/security", page: "nav.security", group: "sec.systemProtection", label: "sec.memIntegrity", desc: "sec.memIntegrityDesc" },
  { path: "/security", page: "nav.security", group: "sec.privacy", label: "sec.telemetryTitle" },
  { path: "/security", page: "nav.security", group: "sec.network", label: "sec.dns", desc: "sec.dnsDesc" },
  { path: "/security", page: "nav.security", group: "sec.network", label: "sec.dnsAdapter" },
  { path: "/security", page: "nav.security", group: "sec.network", label: "sec.dnsPrimary" },
  { path: "/security", page: "nav.security", group: "sec.network", label: "sec.dnsSecondary" },
  { path: "/security", page: "nav.security", group: "sec.ifeoTitle", label: "sec.ifeoDesc" },
  { path: "/updates", page: "nav.updates", group: "upd.components", label: "upd.certs", desc: "upd.certsDesc" },
  { path: "/updates", page: "nav.updates", group: "upd.components", label: "upd.kgl", desc: "upd.kglDesc" },
  { path: "/updates", page: "nav.updates", group: "upd.windowsUpdate", label: "upd.pause", desc: "upd.pauseDesc" },
  { path: "/updates", page: "nav.updates", group: "upd.windowsUpdate", label: "upd.hidePage", desc: "upd.hidePageDesc" },
  { path: "/updates", page: "nav.updates", group: "upd.windowsUpdate", label: "upd.drivers", desc: "upd.driversDesc" },
  { path: "/updates", page: "nav.updates", group: "upd.windowsUpdate", label: "upd.channel", desc: "upd.channelDesc" },
  { path: "/utilities", page: "nav.utilities", group: "util.powerManagement", label: "util.hibernate", desc: "util.hibernateDesc" },
  { path: "/utilities", page: "nav.utilities", group: "util.powerManagement", label: "util.fastStartup", desc: "util.fastStartupDesc" },
  { path: "/utilities", page: "nav.utilities", group: "util.apps", label: "util.photoViewer", desc: "util.photoViewerDesc" },
  { path: "/utilities", page: "nav.utilities", group: "util.apps", label: "util.edgeLabel" },
  { path: "/utilities", page: "nav.utilities", group: "util.apps", label: "util.currentVersion" },
]
