<template>
  <div class="cloudflare-email mt-5">
    <h2 class="title is-5">{{ $t('settings.cloudflare.name') }}</h2>
    <p class="has-text-grey is-size-7 mb-4">{{ $t('settings.cloudflare.help') }}</p>

    <div class="block box">
      <div class="columns">
        <div class="column is-2">
          <b-field :label="$t('globals.buttons.enabled')">
            <b-switch v-model="cf.enabled" name="cloudflare_enabled" :native-value="true"
              data-cy="btn-enable-cloudflare" />
          </b-field>
        </div>

        <div class="column" :class="{ disabled: !cf.enabled }">
          <div class="columns">
            <div class="column is-6">
              <b-field :label="$t('settings.cloudflare.accountID')" label-position="on-border"
                :message="$t('settings.cloudflare.accountIDHelp')">
                <b-input v-model="cf.account_id" name="cloudflare_account_id"
                  placeholder="0123456789abcdef0123456789abcdef" :maxlength="64" />
              </b-field>
            </div>
            <div class="column is-6">
              <b-field :label="$t('settings.cloudflare.apiToken')" label-position="on-border"
                :message="$t('settings.cloudflare.apiTokenHelp')">
                <b-input v-model="cf.api_token" name="cloudflare_api_token" type="password"
                  class="cloudflare-api-token"
                  :placeholder="$t('globals.messages.passwordChange')" :maxlength="200" />
              </b-field>
            </div>
          </div>

          <div class="columns">
            <div class="column is-4">
              <b-field :label="$t('settings.messengers.maxConns')" label-position="on-border"
                :message="$t('settings.messengers.maxConnsHelp')">
                <b-numberinput v-model="cf.max_conns" name="cloudflare_max_conns" type="is-light"
                  controls-position="compact" placeholder="10" min="1" max="65535" />
              </b-field>
            </div>
            <div class="column is-4">
              <b-field :label="$t('settings.messengers.retries')" label-position="on-border"
                :message="$t('settings.messengers.retriesHelp')">
                <b-numberinput v-model="cf.max_msg_retries" name="cloudflare_max_msg_retries"
                  type="is-light" controls-position="compact" placeholder="2" min="1" max="1000" />
              </b-field>
            </div>
            <div class="column is-4">
              <b-field :label="$t('settings.messengers.timeout')" label-position="on-border"
                :message="$t('settings.messengers.timeoutHelp')">
                <b-input v-model="cf.timeout" name="cloudflare_timeout" placeholder="10s"
                  :pattern="regDuration" :maxlength="10" />
              </b-field>
            </div>
          </div>

          <hr />

          <form @submit.prevent="doTest">
            <div>
              <div v-if="showTest">
                <div class="columns">
                  <div class="column is-6">
                    <b-field :label="$t('settings.cloudflare.fromEmail')" label-position="on-border"
                      :message="$t('settings.cloudflare.fromEmailHelp')">
                      <b-input v-model="testFrom" name="cloudflare_test_from"
                        placeholder="hello@yourdomain.com" class="test-from-cloudflare"
                        :maxlength="300" required />
                    </b-field>
                  </div>
                  <div class="column is-6">
                    <b-field :label="$t('settings.smtp.toEmail')" label-position="on-border" grouped>
                      <b-input v-model="testEmail" name="cloudflare_test_email" type="email"
                        :placeholder="$t('subscribers.email')" class="test-email-cloudflare"
                        :maxlength="300" required />
                      <p class="controls">
                        <b-button native-type="submit" type="is-success" :disabled="!isTestEnabled"
                          icon-left="email-outline" data-cy="btn-cloudflare-test">
                          {{ $t('settings.smtp.sendTest') }}
                        </b-button>
                      </p>
                    </b-field>
                  </div>
                </div>
              </div>
              <div v-else>
                <a href="#" class="is-primary" @click.prevent="showTestForm"
                  :class="{ 'has-text-grey': !isTestEnabled }">
                  <b-icon icon="rocket-launch-outline" /> {{ $t('settings.cloudflare.testConnection') }}
                </a>
              </div>
              <div v-if="errMsg">
                <b-field class="mt-4" type="is-danger">
                  <b-input v-model="errMsg" type="textarea" custom-class="has-text-danger is-size-6" readonly />
                </b-field>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import { regDuration } from '../../constants';

const defaultCloudflare = () => ({
  enabled: false,
  account_id: '',
  api_token: '',
  max_conns: 10,
  timeout: '10s',
  max_msg_retries: 2,
});

export default Vue.extend({
  props: {
    form: {
      type: Object, default: () => ({}),
    },
  },

  data() {
    if (!this.form.cloudflare_email) {
      this.$set(this.form, 'cloudflare_email', defaultCloudflare());
    }
    return {
      regDuration,
      showTest: false,
      testFrom: '',
      testEmail: '',
      errMsg: '',
    };
  },

  computed: {
    ...mapState(['settings']),
    cf() {
      return this.form.cloudflare_email;
    },
    isTestEnabled() {
      if (!this.cf.enabled || !this.cf.account_id) {
        return false;
      }
      if (!this.cf.api_token || this.cf.api_token.includes('•')) {
        return false;
      }
      return true;
    },
  },

  methods: {
    showTestForm() {
      if (!this.isTestEnabled) {
        this.$utils.toast(this.$t('settings.cloudflare.testEnterToken'), 'is-danger');
        this.$nextTick(() => {
          const i = document.querySelector('.cloudflare-api-token input');
          if (i) {
            this.cf.api_token = '';
            i.focus();
            i.select();
          }
        });
        return;
      }
      this.showTest = true;
      this.errMsg = '';
      if (!this.testFrom) {
        this.testFrom = this.form['app.from_email'] || this.settings?.['app.from_email'] || '';
      }
      this.$nextTick(() => {
        document.querySelector('.test-from-cloudflare')?.focus();
      });
    },

    doTest() {
      if (!this.isTestEnabled) {
        this.$utils.toast(this.$t('settings.cloudflare.testEnterToken'), 'is-danger');
        return;
      }
      if (!this.testFrom?.trim()) {
        this.$utils.toast(this.$t('settings.cloudflare.fromEmailHelp'), 'is-danger');
        return;
      }

      this.errMsg = '';
      this.$api.testCloudflare({
        account_id: this.cf.account_id,
        api_token: this.cf.api_token,
        max_conns: this.cf.max_conns,
        timeout: this.cf.timeout,
        from: this.testFrom.trim(),
        email: this.testEmail,
      }).then(() => {
        this.$utils.toast(this.$t('campaigns.testSent'));
      }).catch((err) => {
        if (err.response?.data?.message) {
          this.errMsg = err.response.data.message;
        }
      });
    },
  },
});
</script>
