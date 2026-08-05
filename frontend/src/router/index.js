import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

// The meta.group param is used in App.vue to expand menu group by name.
const routes = [
  {
    path: '/404',
    name: '404_page',
    meta: { title: '404' },
    component: () => import('../views/404.vue'),
  },
  {
    path: '/',
    name: 'dashboard',
    meta: { title: '' },
    component: () => import('../views/Dashboard.vue'),
  },
  {
    path: '/lists',
    name: 'lists',
    meta: { title: 'globals.terms.lists', group: 'crm' },
    component: () => import('../views/Lists.vue'),
  },
  {
    path: '/lists/forms',
    redirect: { name: 'forms' },
  },
  {
    path: '/lists/:id',
    redirect: (to) => `/contacts/lists/${to.params.id}`,
  },
  {
    path: '/contacts',
    name: 'subscribers',
    meta: { title: 'menu.contacts', group: 'crm' },
    component: () => import('../views/Subscribers.vue'),
  },
  {
    path: '/contacts/import',
    name: 'import',
    meta: { title: 'import.title', group: 'crm' },
    component: () => import('../views/Import.vue'),
  },
  {
    path: '/contacts/bounces',
    name: 'bounces',
    meta: { title: 'globals.terms.bounces', group: 'crm' },
    component: () => import('../views/Bounces.vue'),
  },
  {
    path: '/contacts/lists/:listID',
    name: 'subscribers_list',
    meta: { title: 'menu.contacts', group: 'crm' },
    component: () => import('../views/Subscribers.vue'),
  },
  {
    path: '/contacts/:id',
    name: 'subscriber',
    meta: { title: 'menu.contacts', group: 'crm' },
    component: () => import('../views/Subscribers.vue'),
  },
  // Legacy subscriber URLs
  {
    path: '/subscribers',
    redirect: '/contacts',
  },
  {
    path: '/subscribers/import',
    redirect: '/contacts/import',
  },
  {
    path: '/subscribers/bounces',
    redirect: '/contacts/bounces',
  },
  {
    path: '/subscribers/lists/:listID',
    redirect: (to) => `/contacts/lists/${to.params.listID}`,
  },
  {
    path: '/subscribers/:id',
    redirect: (to) => `/contacts/${to.params.id}`,
  },
  {
    path: '/segments',
    name: 'segments',
    meta: { title: 'menu.segments', group: 'crm' },
    component: () => import('../views/Segments.vue'),
  },
  {
    path: '/companies',
    name: 'companies',
    meta: { title: 'menu.companies', group: 'crm' },
    component: () => import('../views/Companies.vue'),
  },
  {
    path: '/tasks',
    name: 'tasks',
    meta: { title: 'menu.tasks', group: 'crm' },
    component: () => import('../views/Tasks.vue'),
  },
  {
    path: '/campaigns',
    name: 'campaigns',
    meta: { title: 'globals.terms.campaigns', group: 'marketing' },
    component: () => import('../views/Campaigns.vue'),
  },
  {
    path: '/campaigns/media',
    name: 'media',
    meta: { title: 'globals.terms.media', group: 'marketing' },
    component: () => import('../views/Media.vue'),
  },
  {
    path: '/campaigns/templates',
    redirect: { name: 'templates' },
  },
  {
    path: '/campaigns/templates/:id',
    redirect: (to) => ({ name: 'template', params: { id: to.params.id } }),
  },
  {
    path: '/marketing/forms',
    name: 'forms',
    meta: { title: 'forms.title', group: 'marketing' },
    component: () => import('../views/Forms.vue'),
  },
  {
    path: '/marketing/statistics',
    name: 'statistics',
    meta: { title: 'menu.statistics', group: 'marketing' },
    component: () => import('../views/Statistics.vue'),
  },
  {
    path: '/marketing/templates',
    name: 'templates',
    meta: { title: 'menu.templates', group: 'marketing' },
    component: () => import('../views/Templates.vue'),
  },
  {
    path: '/marketing/templates/:id',
    name: 'template',
    meta: { title: 'menu.templates', group: 'marketing' },
    component: () => import('../views/Template.vue'),
  },
  {
    path: '/campaigns/analytics',
    name: 'campaignAnalytics',
    redirect: { name: 'statistics' },
  },
  {
    path: '/campaigns/:id',
    name: 'campaign',
    meta: { title: 'globals.terms.campaign', group: 'marketing' },
    component: () => import('../views/Campaign.vue'),
  },
  {
    path: '/user/profile',
    name: 'userProfile',
    meta: { title: 'users.profile', group: 'settings' },
    component: () => import('../views/UserProfile.vue'),
  },
  {
    path: '/settings',
    name: 'settings',
    meta: { title: 'globals.terms.settings', group: 'settings' },
    component: () => import('../views/Settings.vue'),
  },
  {
    path: '/settings/logs',
    name: 'logs',
    meta: { title: 'logs.title', group: 'settings' },
    component: () => import('../views/Logs.vue'),
  },
  {
    path: '/users',
    name: 'users',
    meta: { title: 'globals.terms.users', group: 'users' },
    component: () => import('../views/Users.vue'),
  },
  {
    path: '/users/roles/users',
    name: 'userRoles',
    meta: { title: 'users.userRoles', group: 'users' },
    component: () => import('../views/Roles.vue'),
  },
  {
    path: '/users/roles/lists',
    name: 'listRoles',
    meta: { title: 'users.listRoles', group: 'users' },
    component: () => import('../views/Roles.vue'),
  },
  {
    path: '/settings/maintenance',
    name: 'maintenance',
    meta: { title: 'maintenance.title', group: 'settings' },
    component: () => import('../views/Maintenance.vue'),
  },
];

const router = new VueRouter({
  mode: 'history',
  base: import.meta.env.BASE_URL,
  routes,

  scrollBehavior(to) {
    if (to.hash) {
      return { selector: to.hash };
    }
    return { x: 0, y: 0 };
  },
});

export default router;
