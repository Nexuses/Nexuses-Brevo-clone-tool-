<template>
  <nav class="crm-subnav" aria-label="CRM">
    <router-link
      v-for="tab in tabs"
      :key="tab.name"
      :to="{ name: tab.name }"
      class="crm-subnav__tab"
      :class="{ 'is-active': isActive(tab) }"
    >
      <span class="crm-subnav__icon" aria-hidden="true">
        <svg v-if="tab.icon === 'contacts'" width="16" height="16" viewBox="0 0 24 24" fill="none">
          <circle cx="9" cy="8" r="3" stroke="currentColor" stroke-width="1.7" />
          <path d="M4.5 18.5c.8-2.6 2.8-4 4.5-4s3.7 1.4 4.5 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
          <circle cx="16.5" cy="8.5" r="2.3" stroke="currentColor" stroke-width="1.7" />
          <path d="M15.2 14.2c1.4.3 2.7 1.3 3.3 3.3" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
        </svg>
        <svg v-else-if="tab.icon === 'lists'" width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M8 7h11.5M8 12h11.5M8 17h11.5" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
          <circle cx="5" cy="7" r="1.2" fill="currentColor" />
          <circle cx="5" cy="12" r="1.2" fill="currentColor" />
          <circle cx="5" cy="17" r="1.2" fill="currentColor" />
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M9.5 11.5 5 16l1.5 1.5 4.5-4.5" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M12 5.5h7.5v13H12" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" />
          <path d="M14.5 9.5h3M14.5 13h3M14.5 16.5h2" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
        </svg>
      </span>
      <span class="crm-subnav__label">{{ tab.label }}</span>
    </router-link>
  </nav>
</template>

<script>
export default {
  name: 'CrmSubnav',

  computed: {
    tabs() {
      return [
        {
          name: 'subscribers',
          label: 'Contacts',
          match: ['subscribers', 'subscriber', 'subscribers_list'],
          icon: 'contacts',
        },
        {
          name: 'lists',
          label: 'Lists',
          match: ['lists', 'list'],
          icon: 'lists',
        },
        {
          name: 'tasks',
          label: 'Tasks',
          match: ['tasks'],
          icon: 'tasks',
        },
      ];
    },
  },

  methods: {
    isActive(tab) {
      return tab.match.includes(this.$route.name);
    },
  },
};
</script>
