<template>
  <nav class="crm-subnav" aria-label="CRM">
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
  name: 'CrmSubnav',

  computed: {
    tabs() {
      return [
        { name: 'subscribers', label: 'Contacts', match: ['subscribers', 'subscriber', 'subscribers_list'] },
        { name: 'lists', label: 'Lists', match: ['lists', 'list'] },
        { name: 'segments', label: 'Segments', match: ['segments'] },
        { name: 'companies', label: 'Companies', match: ['companies'] },
      ];
    },
  },

  methods: {
    isActive(tab) {
      const routeName = this.$route.name;
      return tab.match.includes(routeName);
    },
  },
};
</script>
