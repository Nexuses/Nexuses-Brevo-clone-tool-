<template>
  <section class="domain-add bv-page">
    <ol class="sdd-stepper" aria-label="Add domain steps">
      <li
        v-for="(label, i) in steps"
        :key="label"
        class="sdd-stepper__item"
        :class="{
          'is-active': step === i + 1,
          'is-done': step > i + 1,
        }"
      >
        <span class="sdd-stepper__num">{{ i + 1 }}</span>
        <span class="sdd-stepper__label">{{ label }}</span>
      </li>
    </ol>

    <form class="domain-add__form" @submit.prevent="onContinue">
      <!-- Step 1: Domain -->
      <div v-if="step === 1" class="domain-add__step">
        <h1 class="domain-add__title">Add your sending domain</h1>
        <label class="sender-add__field">
          <span class="sender-add__label">Domain name <em>*</em></span>
          <p class="sender-add__hint">
            Enter the root domain you own (e.g. eguardian.in or yourcompany.com).
            You can choose a branded tracking subdomain now or set it up later when
            you authenticate the domain.
          </p>
          <input
            v-model="domain"
            class="sender-add__input"
            type="text"
            name="domain"
            maxlength="255"
            placeholder="yourcompany.com"
            autocomplete="off"
            spellcheck="false"
            data-cy="domain"
            required
          >
        </label>
      </div>

      <!-- Step 2: Branded subdomain (optional) -->
      <div v-else-if="step === 2" class="domain-add__step">
        <h1 class="domain-add__title">Choose a branded subdomain <span class="domain-add__optional">(optional)</span></h1>
        <p class="sender-add__hint domain-add__intro">
          Clicks and opens can use a subdomain such as <code>emailtrack</code> so
          tracking links appear under your brand. You can skip this and configure
          the tracking host when you authenticate the domain.
        </p>
        <button
          type="button"
          class="domain-add__skip-link"
          data-cy="btn-skip-subdomain"
          :disabled="isBusy"
          @click="onSkipSubdomain"
        >
          Skip — add <strong>{{ parentDomain }}</strong> only and set tracking host later
        </button>
        <label class="sender-add__field">
          <span class="sender-add__label">Branded subdomain</span>
          <p class="sender-add__hint">
            Optional. Use a short host such as <code>click</code> or <code>emailtrack</code>.
            Do not include your domain name here.
          </p>
          <div class="domain-add__combo">
            <input
              v-model="subdomain"
              class="sender-add__input"
              type="text"
              name="subdomain"
              maxlength="63"
              placeholder="emailtrack (optional)"
              autocomplete="off"
              spellcheck="false"
              data-cy="subdomain"
            >
            <span class="domain-add__combo-suffix">.{{ parentDomain }}</span>
          </div>
        </label>
        <p v-if="hasSubdomain" class="domain-add__preview-url">
          Tracking links will use <strong>{{ fullDomain }}</strong>
        </p>
      </div>

      <!-- Step 3: Setup method -->
      <div v-else-if="step === 3" class="domain-add__step">
        <h1 class="domain-add__title">Choose a setup method</h1>
        <p class="sender-add__hint domain-add__intro">
          Add a CNAME record at your DNS provider so we can verify
          <strong>{{ trackingHostLabel }}</strong>.
        </p>
        <label class="domain-add__method is-selected">
          <input type="radio" name="setup" value="manual" checked>
          <span>
            <strong>Add DNS records yourself</strong>
            <em>Copy the CNAME into your domain registrar or DNS panel.</em>
          </span>
        </label>
      </div>

      <!-- Step 4: Records -->
      <div v-else class="domain-add__step">
        <h1 class="domain-add__title">Add these DNS records</h1>
        <p class="sender-add__hint domain-add__intro">
          Create this CNAME at the DNS provider for <strong>{{ parentDomain }}</strong>,
          then continue. DNS can take a few minutes to 24–48 hours to propagate.
        </p>
        <div class="domain-add__records">
          <table>
            <thead>
              <tr>
                <th>Type</th>
                <th>Name</th>
                <th>Value</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>CNAME</td>
                <td><copy-text :text="dnsName" /></td>
                <td><copy-text :text="cnameTarget || '—'" /></td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-if="!cnameTarget" class="domain-add__warn">
          Set a public tracking URL in Settings → General before this CNAME can verify.
        </p>
      </div>

      <footer class="sdd-wizard__footer">
        <button type="button" class="sdd-wizard__cancel" @click="onCancel">Cancel</button>
        <div class="sdd-wizard__nav">
          <button
            type="button"
            class="sdd-wizard__prev"
            :disabled="step === 1"
            @click="onPrevious"
          >
            Previous
          </button>
          <button
            type="submit"
            class="sdd-wizard__primary"
            :data-cy="step === 4 || (step === 2 && !hasSubdomain) ? 'btn-add' : 'btn-continue'"
            :disabled="!canContinue || isBusy"
          >
            {{ primaryLabel }}
          </button>
        </div>
      </footer>
    </form>
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import CopyText from '../components/CopyText.vue';

export default Vue.extend({
  components: {
    CopyText,
  },

  data() {
    return {
      step: 1,
      domain: '',
      subdomain: '',
      skipSubdomain: false,
      trackingUrl: '',
      rootUrl: '',
      steps: ['Domain', 'Branded subdomain (optional)', 'Setup method', 'Records'],
    };
  },

  computed: {
    ...mapState(['loading']),

    isBusy() {
      return this.loading.trackingDomains;
    },

    parentDomain() {
      return this.normalizeHost(this.domain);
    },

    hasSubdomain() {
      const label = this.normalizeLabel(this.subdomain);
      return label !== '' && this.isValidLabel(label);
    },

    fullDomain() {
      const prefix = this.normalizeLabel(this.subdomain);
      const parent = this.parentDomain;
      if (!prefix || !parent) return '';
      return `${prefix}.${parent}`;
    },

    trackingHostLabel() {
      if (this.skipSubdomain || !this.isValidLabel(this.normalizeLabel(this.subdomain))) {
        return this.parentDomain;
      }
      return this.fullDomain;
    },

    dnsName() {
      if (this.skipSubdomain || !this.isValidLabel(this.normalizeLabel(this.subdomain))) {
        return 'Set when authenticating';
      }
      return this.normalizeLabel(this.subdomain);
    },

    cnameTarget() {
      return this.hostFromUrl(this.trackingUrl) || this.hostFromUrl(this.rootUrl);
    },

    primaryLabel() {
      if (this.step === 4) return 'Add domain';
      if (this.step === 2 && !this.hasSubdomain) return 'Skip & add domain';
      return 'Continue';
    },

    canContinue() {
      if (this.step === 1) {
        return this.isValidDomain(this.parentDomain);
      }
      if (this.step === 2) {
        if (!this.hasSubdomain) return true;
        return this.isValidDomain(this.fullDomain);
      }
      return true;
    },
  },

  methods: {
    normalizeHost(raw) {
      return String(raw || '')
        .trim()
        .toLowerCase()
        .replace(/^https?:\/\//, '')
        .replace(/^www\./, '')
        .replace(/\/.*$/, '')
        .replace(/\.$/, '');
    },

    normalizeLabel(raw) {
      return String(raw || '')
        .trim()
        .toLowerCase()
        .replace(/\.$/, '')
        .split('.')[0];
    },

    hostFromUrl(raw) {
      const value = String(raw || '').trim();
      if (!value) return '';
      const withScheme = value.indexOf('://') > -1 ? value : `https://${value}`;
      const hostPort = withScheme.replace(/^https?:\/\//i, '').split('/')[0];
      return hostPort.split(':')[0].replace(/\.$/, '').toLowerCase();
    },

    isValidLabel(label) {
      return /^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$/.test(label);
    },

    isValidDomain(domain) {
      return /^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+$/.test(domain);
    },

    onCancel() {
      this.$router.push({ name: 'trackingDomains', query: { tab: 'domains' } });
    },

    onPrevious() {
      if (this.step > 1) {
        this.skipSubdomain = false;
        this.step -= 1;
      }
    },

    onSkipSubdomain() {
      this.skipSubdomain = true;
      this.subdomain = '';
      this.createDomain();
    },

    createDomain() {
      const payload = { domain: this.parentDomain };
      if (this.hasSubdomain) {
        payload.tracking_host = this.fullDomain;
      }
      this.$api.createTrackingDomain(payload).then(() => {
        const msg = this.skipSubdomain
          ? 'Domain added. Choose a tracking subdomain when you authenticate it.'
          : 'Domain added. Add the DNS record, then verify it.';
        this.$utils.toast(msg);
        this.$router.push({ name: 'trackingDomains', query: { tab: 'domains' } });
      });
    },

    onContinue() {
      if (!this.canContinue) return;
      if (this.step === 2 && !this.hasSubdomain) {
        this.skipSubdomain = true;
        this.createDomain();
        return;
      }
      if (this.step < 4) {
        this.step += 1;
        return;
      }
      this.createDomain();
    },
  },

  mounted() {
    this.$api.getSettings().then((data) => {
      this.trackingUrl = (data && data['app.tracking_url']) || '';
      this.rootUrl = (data && data['app.root_url']) || '';
    });
  },
});
</script>
