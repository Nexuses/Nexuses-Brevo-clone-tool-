<template>
  <nav class="crm-subnav" aria-label="Marketing">
    <router-link
      v-for="tab in tabs"
      :key="tab.name"
      :to="{ name: tab.name }"
      class="crm-subnav__tab"
      :class="{ 'is-active': isActive(tab) }"
    >
      {{ tab.label }}
    </router-link>
  </nav>
</template>

<script>
export default {
  name: 'MarketingSubnav',

  computed: {
    tabs() {
      return [
        { name: 'campaigns', label: 'Campaigns', match: ['campaigns', 'campaign'] },
        { name: 'forms', label: 'Forms', match: ['forms'] },
        { name: 'statistics', label: 'Statistics', match: ['statistics'] },
        { name: 'templates', label: 'Templates', match: ['templates', 'template'] },
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
