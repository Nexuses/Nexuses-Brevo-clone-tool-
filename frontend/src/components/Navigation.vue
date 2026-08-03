<template>
  <b-menu-list>
    <b-menu-item :to="{ name: 'dashboard' }" tag="router-link" :active="activeItem.dashboard"
      icon="home-outline" :label="$t('menu.dashboard')" /><!-- home -->

    <b-menu-item :expanded="activeGroup.crm" :active="activeGroup.crm" data-cy="crm"
      @update:active="(state) => toggleGroup('crm', state)" icon="account-multiple"
      :label="$t('menu.crm')">
      <b-menu-item v-if="$can('subscribers:get_all', 'subscribers:get')" :to="{ name: 'subscribers' }" tag="router-link"
        :active="activeItem.subscribers || activeItem.subscriber || activeItem.subscribers_list"
        data-cy="crm-contacts" icon="account-multiple"
        :label="$t('menu.contacts')" />
      <b-menu-item :to="{ name: 'lists' }" tag="router-link"
        :active="activeItem.lists || activeItem.list" data-cy="crm-lists"
        icon="format-list-bulleted-square" :label="$t('menu.lists')" />
      <b-menu-item :to="{ name: 'segments' }" tag="router-link" :active="activeItem.segments"
        data-cy="crm-segments" icon="tag-outline" :label="$t('menu.segments')" />
      <b-menu-item :to="{ name: 'companies' }" tag="router-link" :active="activeItem.companies"
        data-cy="crm-companies" icon="newspaper-variant-outline" :label="$t('menu.companies')" />
      <b-menu-item v-if="$can('subscribers:import')" :to="{ name: 'import' }" tag="router-link"
        :active="activeItem.import" data-cy="import" icon="file-upload-outline" :label="$t('menu.import')" />
      <b-menu-item v-if="$can('bounces:get')" :to="{ name: 'bounces' }" tag="router-link" :active="activeItem.bounces"
        data-cy="bounces" icon="email-bounce" :label="$t('globals.terms.bounces')" />
    </b-menu-item><!-- crm -->

    <b-menu-item v-if="$can('campaigns:*', 'templates:*')" :expanded="activeGroup.marketing" :active="activeGroup.marketing"
      data-cy="marketing" @update:active="(state) => toggleGroup('marketing', state)" icon="rocket-launch-outline"
      :label="$t('menu.marketing')">
      <b-menu-item v-if="$can('campaigns:get')" :to="{ name: 'campaigns' }" tag="router-link"
        :active="activeItem.campaigns || activeItem.campaign" data-cy="all-campaigns" icon="rocket-launch-outline"
        :label="$t('menu.campaigns')" />
      <b-menu-item :to="{ name: 'forms' }" tag="router-link" :active="activeItem.forms" class="forms"
        icon="newspaper-variant-outline" :label="$t('menu.forms')" />
      <b-menu-item v-if="$can('campaigns:get_analytics')" :to="{ name: 'statistics' }" tag="router-link"
        :active="activeItem.statistics || activeItem.campaignAnalytics" data-cy="statistics" icon="chart-bar"
        :label="$t('menu.statistics')" />
      <b-menu-item v-if="$can('templates:get')" :to="{ name: 'templates' }" tag="router-link"
        :active="activeItem.templates || activeItem.template" data-cy="templates" icon="file-image-outline"
        :label="$t('menu.templates')" />
      <b-menu-item v-if="$can('campaigns:manage')" :to="{ name: 'campaign', params: { id: 'new' } }" tag="router-link"
        data-cy="new-campaign" icon="plus" :label="$t('menu.newCampaign')" />
      <b-menu-item v-if="$can('media:*')" :to="{ name: 'media' }" tag="router-link" :active="activeItem.media"
        data-cy="media" icon="image-outline" :label="$t('menu.media')" />
    </b-menu-item><!-- marketing -->

    <b-menu-item v-if="$can('users:*', 'roles:*')" :expanded="activeGroup.users" :active="activeGroup.users"
      data-cy="users" @update:active="(state) => toggleGroup('users', state)" icon="account-multiple"
      :label="$t('globals.terms.users')">
      <b-menu-item v-if="$can('users:get')" :to="{ name: 'users' }" tag="router-link" :active="activeItem.users"
        data-cy="users" icon="account-multiple" :label="$t('globals.terms.users')" />
      <b-menu-item v-if="$can('roles:get')" :to="{ name: 'userRoles' }" tag="router-link" :active="activeItem.userRoles"
        data-cy="userRoles" icon="newspaper-variant-outline" :label="$t('users.userRoles')" />
      <b-menu-item v-if="$can('roles:get')" :to="{ name: 'listRoles' }" tag="router-link" :active="activeItem.listRoles"
        data-cy="listRoles" icon="format-list-bulleted-square" :label="$t('users.listRoles')" />
    </b-menu-item><!-- users -->

    <b-menu-item v-if="$can('settings:*')" :expanded="activeGroup.settings" :active="activeGroup.settings"
      data-cy="settings" @update:active="(state) => toggleGroup('settings', state)" icon="cog-outline"
      :label="$t('menu.settings')">
      <b-menu-item v-if="$can('settings:get')" :to="{ name: 'settings' }" tag="router-link"
        :active="activeItem.settings" data-cy="all-settings" icon="cog-outline" :label="$t('menu.settings')" />
      <b-menu-item v-if="$can('settings:maintain')" :to="{ name: 'maintenance' }" tag="router-link"
        :active="activeItem.maintenance" data-cy="maintenance" icon="wrench-outline" :label="$t('menu.maintenance')" />
      <b-menu-item v-if="$can('settings:get')" :to="{ name: 'logs' }" tag="router-link" :active="activeItem.logs"
        data-cy="logs" icon="format-list-bulleted-square" :label="$t('menu.logs')" />
    </b-menu-item><!-- settings -->
  </b-menu-list>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'Navigation',

  props: {
    activeItem: { type: Object, default: () => { } },
    activeGroup: { type: Object, default: () => { } },
    isMobile: Boolean,
  },

  methods: {
    toggleGroup(group, state) {
      this.$emit('toggleGroup', group, state);
    },

    doLogout() {
      this.$emit('doLogout');
    },
  },

  computed: {
    ...mapState(['profile']),
  },

  mounted() {
    if (this.isMobile) {
      document.querySelectorAll('.navbar li a[href]').forEach((e) => {
        e.onclick = () => {
          document.querySelector('.navbar-burger').click();
        };
      });
    }
  },
};

</script>
