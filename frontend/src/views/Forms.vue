<template>
  <section class="forms forms-brevo bv-page content relative">
    <marketing-subnav />

    <header class="forms-brevo__header">
      <div>
        <h1 class="forms-brevo__title">
          {{ $t('forms.title') }}
        </h1>
        <p class="forms-brevo__lead">
          Build embeddable subscription forms for your public lists.
        </p>
      </div>
    </header>

    <b-loading v-if="loading.lists" :active="loading.lists" :is-full-page="false" />
    <p v-else-if="publicLists.length === 0" class="bv-page__card">
      {{ $t('forms.noPublicLists') }}
    </p>
    <div class="columns" v-else-if="publicLists.length > 0">
      <div class="column is-4">
        <h4>{{ $t('forms.publicLists') }}</h4>
        <p class="has-text-grey">{{ $t('forms.selectHelp') }}</p>

        <b-loading :active="loading.lists" :is-full-page="false" />
        <ul class="no" data-cy="lists">
          <li v-for="(l, i) in publicLists" :key="l.id">
            <b-checkbox v-model="checked" :native-value="i">
              {{ l.name }}
            </b-checkbox>
          </li>
        </ul>

        <template v-if="publicSubEnabled">
          <hr />
          <h4>{{ $t('forms.publicSubPage') }}</h4>
          <p>
            <a :href="`${serverConfig.root_url}/subscription/form`" target="_blank" rel="noopener noreferer"
              data-cy="url" class="bv-link">
              {{ serverConfig.root_url }}/subscription/form
            </a>
          </p>
        </template>
      </div>
      <div class="column" data-cy="form">
        <h4>{{ $t('forms.formHTML') }}</h4>
        <p>
          {{ $t('forms.formHTMLHelp') }}
        </p>

        <code-editor lang="html" v-if="checked.length > 0" v-model="html" disabled />
      </div>
    </div><!-- columns -->
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import CodeEditor from '../components/CodeEditor.vue';
import MarketingSubnav from '../components/MarketingSubnav.vue';

export default Vue.extend({
  name: 'ListForm',

  components: {
    'code-editor': CodeEditor,
    MarketingSubnav,
  },

  data() {
    return {
      checked: [],
      html: '',
    };
  },

  methods: {
    renderHTML() {
      let h = `<form method="post" action="${this.serverConfig.root_url}/subscription/form" class="listmonk-form">\n`
        + '  <div>\n'
        + `    <h3>${this.$t('public.sub')}</h3>\n`
        + '    <input type="hidden" name="nonce" />\n\n'
        + `    <p><input type="email" name="email" required placeholder="${this.$t('subscribers.email')}" /></p>\n`
        + `    <p><input type="text" name="name" placeholder="${this.$t('public.subName')}" /></p>\n\n`;

      this.checked.forEach((i) => {
        const l = this.publicLists[parseInt(i, 10)];

        h += '    <p>\n'
          + `      <input id="${l.uuid.substr(0, 5)}" type="checkbox" name="l" checked value="${l.uuid}" />\n`
          + `      <label for="${l.uuid.substr(0, 5)}">${l.name}</label>\n`;

        if (l.description) {
          h += '      <br />\n'
            + `      <span>${l.description}</span>\n`;
        }

        h += '    </p>\n';
      });

      // Captcha?
      const pubSub = this.serverConfig && this.serverConfig.public_subscription;
      if (pubSub && pubSub.captcha_enabled) {
        if (pubSub.captcha_provider === 'altcha') {
          h += '\n'
            + `    <altcha-widget challengeurl="${this.serverConfig.root_url}/api/public/captcha/altcha"></altcha-widget>\n`
            + `    <${'script'} type="module" src="${this.serverConfig.root_url}/public/static/altcha.umd.js" async defer></${'script'}>\n`;
        } else if (pubSub.captcha_provider === 'hcaptcha') {
          h += '\n'
            + `    <div class="h-captcha" data-sitekey="${pubSub.captcha_key}"></div>\n`
            + `    <${'script'} src="https://js.hcaptcha.com/1/api.js" async defer></${'script'}>\n`;
        }
      }

      h += '\n'
        + `    <input type="submit" value="${this.$t('public.sub')} " />\n`
        + '  </div>\n'
        + '</form>';

      this.html = h;
    },
  },

  computed: {
    ...mapState(['loading', 'lists', 'serverConfig']),

    publicLists() {
      if (!this.lists || !this.lists.results) {
        return [];
      }
      return this.lists.results.filter((l) => l.type === 'public');
    },

    publicSubEnabled() {
      return !!(this.serverConfig && this.serverConfig.public_subscription
        && this.serverConfig.public_subscription.enabled);
    },
  },

  watch: {
    checked() {
      this.renderHTML();
    },
  },

  mounted() {
    this.$api.getLists({ minimal: true, per_page: 'all', status: 'active' });
  },
});
</script>
