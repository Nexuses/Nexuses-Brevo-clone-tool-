<template>
  <section class="maintenance maintenance-brevo wrap bv-page">
    <header class="maintenance-brevo__header">
      <h1 class="maintenance-brevo__title">
        {{ $t('maintenance.title') }}
      </h1>
      <p class="maintenance-brevo__lead">
        {{ $t('maintenance.help') }}
      </p>
    </header>

    <div class="bv-panel">
      <h2 class="bv-panel__title">
        {{ $t('globals.terms.subscribers') }}
      </h2>
      <div class="columns">
        <div class="column is-4">
          <b-field label="Data" :message="$t('maintenance.orphanHelp')">
            <b-select v-model="subscriberType" expanded>
              <option value="orphan">
                {{ $t('dashboard.orphanSubs') }}
              </option>
              <option value="blocklisted">
                {{ $t('subscribers.status.blocklisted') }}
              </option>
            </b-select>
          </b-field>
        </div>
        <div class="column is-5" />
        <div class="column">
          <b-field>
            <b-button type="is-dark" class="btn-new" :loading="loading.maintenance" @click="deleteSubscribers" expanded>
              {{ $t('globals.buttons.delete') }}
            </b-button>
          </b-field>
        </div>
      </div>
    </div><!-- subscribers -->

    <div class="bv-panel">
      <h2 class="bv-panel__title">
        {{ $tc('globals.terms.subscriptions', 2) }}
      </h2>
      <div class="columns">
        <div class="column is-4">
          <b-field label="Data">
            <b-select v-model="subscriptionType" expanded>
              <option value="optin">
                {{ $t('maintenance.maintenance.unconfirmedOptins') }}
              </option>
            </b-select>
          </b-field>
        </div>
        <div class="column is-4">
          <b-field :label="$t('maintenance.olderThan')">
            <b-datepicker v-model="subscriptionDate" required expanded icon="calendar-clock"
              :date-formatter="formatDateTime" />
          </b-field>
        </div>
        <div class="column is-1" />
        <div class="column">
          <b-field>
            <b-button type="is-dark" class="btn-new" :loading="loading.maintenance" @click="deleteSubscriptions" expanded>
              {{ $t('globals.buttons.delete') }}
            </b-button>
          </b-field>
        </div>
      </div>
    </div><!-- subscriptions -->

    <div class="bv-panel">
      <h2 class="bv-panel__title">
        {{ $t('globals.terms.analytics') }}
      </h2>
      <div class="columns">
        <div class="column is-4">
          <b-field label="Data">
            <b-select v-model="analyticsType" expanded>
              <option selected value="all">
                {{ $t('globals.terms.all') }}
              </option>
              <option value="views">
                {{ $t('dashboard.campaignViews') }}
              </option>
              <option value="clicks">
                {{ $t('dashboard.linkClicks') }}
              </option>
            </b-select>
          </b-field>
        </div>
        <div class="column is-4">
          <b-field :label="$t('maintenance.olderThan')">
            <b-datepicker v-model="analyticsDate" required expanded icon="calendar-clock"
              :date-formatter="formatDateTime" />
          </b-field>
        </div>
        <div class="column is-1" />
        <div class="column">
          <b-field>
            <b-button expanded type="is-dark" class="btn-new" :loading="loading.maintenance" @click="deleteAnalytics">
              {{ $t('globals.buttons.delete') }}
            </b-button>
          </b-field>
        </div>
      </div>

      <hr />
      <h3 class="bv-panel__subtitle">
        {{ $t('subscribers.export') }}
      </h3>
      <div class="columns">
        <div class="column is-4">
          <b-field label="Data">
            <b-select v-model="exportType" expanded>
              <option value="views">
                {{ $t('dashboard.campaignViews') }}
              </option>
              <option value="clicks">
                {{ $t('dashboard.linkClicks') }}
              </option>
            </b-select>
          </b-field>
        </div>
        <div class="column is-4">
          <b-field :label="$t('analytics.fromDate')">
            <b-datepicker v-model="exportDate" required expanded icon="calendar-clock"
              :date-formatter="formatDateTime" />
          </b-field>
        </div>
        <div class="column is-1" />
        <div class="column">
          <b-field>
            <b-button expanded type="is-dark" class="btn-new" tag="a" icon-left="download"
              :href="exportURL">
              {{ $t('subscribers.export') }}
            </b-button>
          </b-field>
        </div>
      </div>
    </div><!-- analytics -->

    <form @submit.prevent="onUpdateDBSettings" class="bv-panel">
      <h2 class="bv-panel__title">
        {{ $t('maintenance.database.title') }}
      </h2>
      <h3 class="bv-panel__subtitle">Vacuum</h3>
      <p class="bv-panel__help">
        {{ $t('maintenance.database.vacuumHelp') }}
      </p>
      <div class="columns">
        <div class="column is-2">
          <b-field :label="$t('globals.buttons.enabled')">
            <b-switch v-model="dbSettings.vacuum" />
          </b-field>
        </div>
        <div class="column is-4" :class="{ disabled: !dbSettings.vacuum }">
          <b-field :label="$t('settings.maintenance.cron')">
            <b-input v-model="dbSettings.vacuum_cron_interval" placeholder="0 2 * * *" :disabled="!dbSettings.vacuum"
              pattern="((\*|[0-9,\-\/]+)\s+){4}(\*|[0-9,\-\/]+)" />
          </b-field>
        </div>
        <div class="column is-3" />
        <div class="column is-3">
          <b-button type="is-dark" class="btn-new" native-type="submit" :loading="loading.settings" expanded>
            {{ $t('globals.buttons.save') }}
          </b-button>
        </div>
      </div>
    </form><!-- database -->

    <b-loading :is-full-page="true" v-if="isLoading" active />
  </section>
</template>

<script>
import dayjs from 'dayjs';
import Vue from 'vue';
import { mapState } from 'vuex';

export default Vue.extend({
  components: {
  },

  data() {
    return {
      isLoading: false,
      subscriberType: 'orphan',
      analyticsType: 'all',
      subscriptionType: 'optin',
      analyticsDate: dayjs().subtract(7, 'day').toDate(),
      subscriptionDate: dayjs().subtract(7, 'day').toDate(),
      exportType: 'views',
      exportDate: dayjs().subtract(30, 'day').toDate(),
      dbSettings: {
        vacuum: false,
        vacuum_cron_interval: '0 2 * * *',
      },
    };
  },

  mounted() {
    this.loadDBSettings();
  },

  methods: {
    formatDateTime(s) {
      return dayjs(s).format('YYYY-MM-DD');
    },

    deleteSubscribers() {
      this.$utils.confirm(
        null,
        () => {
          this.$api.deleteGCSubscribers(this.subscriberType).then((data) => {
            this.$utils.toast(this.$t(
              'globals.messages.deletedCount',
              { name: this.$tc('globals.terms.subscribers', 2), num: data.count },
            ));
          });
        },
      );
    },

    deleteSubscriptions() {
      this.$utils.confirm(
        null,
        () => {
          this.$api.deleteGCSubscriptions(this.subscriptionDate).then((data) => {
            this.$utils.toast(this.$t(
              'globals.messages.deletedCount',
              { name: this.$tc('globals.terms.subscriptions', 2), num: data.count },
            ));
          });
        },
      );
    },

    deleteAnalytics() {
      this.$utils.confirm(
        null,
        () => {
          this.$api.deleteGCCampaignAnalytics(this.analyticsType, this.analyticsDate)
            .then(() => {
              this.$utils.toast(this.$t('globals.messages.done'));
            });
        },
      );
    },

    loadDBSettings() {
      this.$api.getSettings().then((data) => {
        if (data['maintenance.db'] !== undefined) {
          this.dbSettings = { ...data['maintenance.db'] };
        }
      });
    },

    async onUpdateDBSettings() {
      this.isLoading = true;
      const data = await this.$api.updateSettingsByKey('maintenance.db', this.dbSettings);
      await this.$root.awaitRestart(data);
      this.isLoading = false;
    },
  },

  computed: {
    ...mapState(['loading']),

    exportURL() {
      const since = encodeURIComponent(dayjs(this.exportDate).toISOString());
      return `/api/maintenance/analytics/${this.exportType}/export?since=${since}`;
    },
  },

});
</script>
