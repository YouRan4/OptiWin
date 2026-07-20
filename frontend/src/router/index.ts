import { createRouter, createWebHashHistory } from 'vue-router'
import HomePage from '../views/HomePage.vue'
import SecurityPage from '../views/SecurityPage.vue'
import PerformancePage from '../views/PerformancePage.vue'
import PersonalizationPage from '../views/PersonalizationPage.vue'
import UtilitiesPage from '../views/UtilitiesPage.vue'
import UpdatesPage from '../views/UpdatesPage.vue'

const routes = [
  { path: '/', name: 'home', component: HomePage },
  { path: '/security', name: 'security', component: SecurityPage },
  { path: '/performance', name: 'performance', component: PerformancePage },
  { path: '/personalization', name: 'personalization', component: PersonalizationPage },
  { path: '/utilities', name: 'utilities', component: UtilitiesPage },
  { path: '/updates', name: 'updates', component: UpdatesPage },
]

export default createRouter({
  history: createWebHashHistory(),
  routes,
})
