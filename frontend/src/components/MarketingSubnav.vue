<template>
  <nav class="crm-subnav" aria-label="Marketing">
    <router-link
      v-for="tab in tabs"
      :key="tab.name"
      :to="{ name: tab.name }"
      class="crm-subnav__tab"
      :class="{ 'is-active': isActive(tab) }"
    >
      <span class="crm-subnav__icon" aria-hidden="true">
        <!-- Campaigns / megaphone -->
        <svg v-if="tab.icon === 'campaigns'" width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M3 11.5v1c0 1.1.9 2 2 2h1l4.5 4.5V5L6 9.5H5c-1.1 0-2 .9-2 2z" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" />
          <path d="M14.5 8.5c1.1.8 1.8 2.1 1.8 3.5s-.7 2.7-1.8 3.5M17.2 6c1.8 1.3 3 3.3 3 5.5s-1.2 4.2-3 5.5" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
        </svg>
        <!-- Forms -->
        <svg v-else-if="tab.icon === 'forms'" width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M7 3.5h7.5L19 8v12.5a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1v-16a1 1 0 0 1 1-1z" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" />
          <path d="M14 3.5V8h4.5M8.5 12h7M8.5 15.5h7M8.5 19h4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        <!-- Statistics -->
        <svg v-else-if="tab.icon === 'statistics'" width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M4.5 19.5h15M7 16.5v-5M12 16.5V7.5M17 16.5v-8" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
        </svg>
        <!-- Templates -->
        <svg v-else-if="tab.icon === 'templates'" width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M4.5 5.5h15v13h-15z" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" />
          <path d="M4.5 9.5h15M10 9.5v9" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
        </svg>
        <!-- Media -->
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M4.5 6.5h15v11h-15z" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" />
          <path d="m8 14.5 2.5-2.5 2 2 3-3.5 2.5 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
          <circle cx="9" cy="9.5" r="1.2" fill="currentColor" />
        </svg>
      </span>
      <span class="crm-subnav__label">{{ tab.label }}</span>
    </router-link>
  </nav>
</template>

<script>
export default {
  name: 'MarketingSubnav',

  computed: {
    tabs() {
      return [
        {
          name: 'campaigns',
          label: 'Campaigns',
          match: ['campaigns', 'campaign'],
          icon: 'campaigns',
        },
        {
          name: 'forms',
          label: 'Forms',
          match: ['forms'],
          icon: 'forms',
        },
        {
          name: 'statistics',
          label: 'Statistics',
          match: ['statistics'],
          icon: 'statistics',
        },
        {
          name: 'templates',
          label: 'Templates',
          match: ['templates', 'template'],
          icon: 'templates',
        },
        {
          name: 'media',
          label: 'Media',
          match: ['media'],
          icon: 'media',
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
