<template>
  <section class="tracking-domains tracking-domains-brevo bv-page">
    <header class="tracking-domains-brevo__header">
      <div class="tracking-domains-brevo__header-main">
        <h1 class="tracking-domains-brevo__title">Senders &amp; Domains</h1>
      </div>
      <button
        v-if="canManage"
        type="button"
        class="tracking-domains-brevo__create"
        data-cy="btn-header-add"
        :disabled="isBusy"
        @click="onHeaderAdd"
      >
        <span class="tracking-domains-brevo__create-plus" aria-hidden="true">+</span>
        {{ activeTab === 'domains' ? 'Add a domain' : 'Add sender' }}
      </button>
    </header>

    <nav class="sdd-tabs" aria-label="Senders and domains">
      <button
        type="button"
        class="sdd-tabs__tab"
        :class="{ 'is-active': activeTab === 'senders' }"
        data-cy="tab-senders"
        @click="setTab('senders')"
      >
        Senders
      </button>
      <button
        type="button"
        class="sdd-tabs__tab"
        :class="{ 'is-active': activeTab === 'domains' }"
        data-cy="tab-domains"
        @click="setTab('domains')"
      >
        Domains
      </button>
    </nav>

    <!-- Senders -->
    <div v-if="activeTab === 'senders'" class="sdd-section" data-cy="senders-panel">
      <h2 class="sdd-section__title">Senders</h2>
      <p class="sdd-section__lead">
        A sender is the name and email address that help recipients recognize your brand.
        To use a sender, it must be saved here. For better email deliverability, authenticate
        the sender's domain on the Domains tab.
      </p>

      <div class="sdd-alert">
        <b-icon icon="warning-empty" />
        <p>
          Google, Yahoo, and Microsoft require authenticated sending domains.
          Add and verify a tracking domain on the Domains tab before sending campaigns.
        </p>
      </div>

      <div class="sdd-search">
        <b-input
          v-model="senderQuery"
          icon="magnify"
          placeholder="Search sender by name or email"
          data-cy="search-sender"
        />
      </div>

      <div v-if="filteredSenders.length === 0 && !loading.settings" class="sdd-empty">
        <empty-placeholder
          label="No senders yet"
          description="Add a sender name and email to use as the From address on campaigns."
        />
      </div>

      <div v-else class="sdd-list">
        <article
          v-for="s in filteredSenders"
          :key="s.email"
          class="sdd-row"
          :class="{ 'is-open': expandedSender === s.email }"
          data-cy="sender-row"
        >
          <header class="sdd-row__head">
            <div class="sdd-row__id">
              <strong>{{ senderLabel(s) }}</strong>
              <span class="sdd-dot sdd-dot--ok">Verified</span>
            </div>
            <div class="sdd-row__actions">
              <button
                v-if="canManage"
                type="button"
                class="sdd-btn"
                data-cy="btn-edit-sender"
                @click="openSenderForm(s)"
              >
                Edit
              </button>
              <button
                type="button"
                class="sdd-icon-btn"
                :aria-expanded="expandedSender === s.email ? 'true' : 'false'"
                aria-label="Toggle sender details"
                @click="toggleSender(s.email)"
              >
                <b-icon :icon="expandedSender === s.email ? 'arrow-up' : 'arrow-down'" />
              </button>
            </div>
          </header>
          <div v-if="expandedSender === s.email" class="sdd-row__body">
            <div class="sdd-meta">
              <div>
                <span class="sdd-meta__lbl">IP address</span>
                <span class="sdd-meta__val">Shared IP</span>
              </div>
              <div>
                <span class="sdd-meta__lbl">Sender domain</span>
                <span class="sdd-meta__val">
                  {{ s.domain || '—' }}
                  <b-icon
                    :icon="domainForEmail(s.email) ? 'check-circle-outline' : 'warning-empty'"
                    size="is-small"
                    :class="domainForEmail(s.email) ? 'has-text-success' : 'has-text-danger'"
                  />
                </span>
              </div>
              <div>
                <span class="sdd-meta__lbl">Tracking</span>
                <span class="sdd-meta__val">
                  {{ domainForEmail(s.email) ? 'Authenticated' : 'Not authenticated' }}
                </span>
              </div>
            </div>
          </div>
        </article>
      </div>
    </div>

    <!-- Domains -->
    <div v-else class="sdd-section" data-cy="domains-panel">
      <h2 class="sdd-section__title">Domains</h2>
      <p class="sdd-section__lead">
        An email domain is the part of your address after the @. Authenticate a tracking
        subdomain such as <code>click.example.com</code> so clicks and opens use your domain.
      </p>
      <p class="sdd-help">
        <b-icon icon="file-find-outline" size="is-small" />
        Add a CNAME record, then authenticate the domain. DNS can take time to propagate.
      </p>

      <div class="sdd-search">
        <b-input
          v-model="domainQuery"
          icon="magnify"
          placeholder="Search domain by name"
          data-cy="search-domain"
        />
      </div>

      <div v-if="filteredDomains.length === 0 && !isBusy" class="sdd-empty">
        <empty-placeholder
          label="No tracking domains yet"
          description="Add a subdomain like click.example.com to send click and open tracking through your own domain."
        />
      </div>

      <div v-else class="sdd-table-wrap">
        <table class="sdd-table" data-cy="domains">
          <thead>
            <tr>
              <th>Domain name</th>
              <th>Authentication status</th>
              <th class="is-actions" />
            </tr>
          </thead>
          <tbody>
            <template v-for="d in filteredDomains">
              <tr :key="d.id" class="sdd-table__row" :data-id="d.id" data-cy="domain-card">
                <td>
                  <strong class="sdd-table__name" data-cy="domain-name">
                    {{ displayDomainName(d) }}
                  </strong>
                  <span
                    v-if="normalizeHost(d.domain) !== baseDomainFor(d)"
                    class="sdd-table__tracking-host"
                    :title="`Tracking host: ${normalizeHost(d.domain)}`"
                  >
                    {{ normalizeHost(d.domain) }}
                  </span>
                </td>
                <td>
                  <span
                    class="sdd-dot"
                    :class="d.status === 'verified' ? 'sdd-dot--ok' : 'sdd-dot--bad'"
                    data-cy="status"
                  >
                    {{ d.status === 'verified' ? 'Authenticated' : 'Not authenticated' }}
                  </span>
                </td>
                <td class="is-actions">
                  <button
                    v-if="canManage"
                    type="button"
                    class="sdd-link"
                    data-cy="btn-verify"
                    @click="onDomainAction(d)"
                  >
                    {{ d.status === 'verified' ? 'View configuration' : 'Authenticate' }}
                  </button>
                  <button
                    v-if="canManage"
                    type="button"
                    class="sdd-icon-btn is-danger"
                    data-cy="btn-delete"
                    :disabled="isBusy"
                    aria-label="Remove domain"
                    @click="onDeleteDomain(d)"
                  >
                    <b-icon icon="trash-can-outline" />
                  </button>
                </td>
              </tr>
              <tr v-if="expandedDomain === d.id" :key="`${d.id}-cfg`" class="sdd-table__config">
                <td colspan="3">
                  <p v-if="d.lastError" class="ctd-card__error" data-cy="last-error">
                    {{ d.lastError }}
                  </p>
                  <div class="ctd-dns">
                    <h3 class="ctd-dns__title">DNS record</h3>
                    <p class="ctd-dns__help">
                      Add this CNAME at your DNS provider for
                      <strong>{{ baseDomainFor(d) }}</strong>, then click
                      <strong>Verify DNS</strong>.
                    </p>
                    <label
                      v-if="canManage && d.status !== 'verified'"
                      class="ctd-dns__host-field"
                    >
                      <span class="ctd-dns__host-label">Tracking hostname</span>
                      <input
                        v-model="trackingHostDraft[d.id]"
                        class="ctd-dns__host-input"
                        type="text"
                        :placeholder="`emailtrack.${baseDomainFor(d)}`"
                        autocomplete="off"
                        spellcheck="false"
                        data-cy="tracking-host"
                        @input="onTrackingHostInput(d)"
                      >
                      <span class="ctd-dns__host-hint">
                        Enter the full tracking host (e.g.
                        <code>emailtrack.{{ baseDomainFor(d) }}</code>).
                        The record name below is the subdomain part only.
                      </span>
                    </label>
                    <dl class="ctd-dns__grid">
                      <div class="ctd-dns__item">
                        <dt>Type</dt>
                        <dd><copy-text :text="d.dnsRecordType" /></dd>
                      </div>
                      <div class="ctd-dns__item">
                        <dt>Name</dt>
                        <dd><copy-text :text="dnsHostLabelFor(d)" /></dd>
                      </div>
                      <div class="ctd-dns__item is-value">
                        <dt>Value</dt>
                        <dd><copy-text :text="d.dnsRecordValue" /></dd>
                      </div>
                    </dl>
                    <p class="ctd-dns__note">
                      DNS changes can take anywhere from a few minutes to 24&ndash;48 hours.
                    </p>
                    <div v-if="canManage" class="ctd-card__actions">
                      <b-button
                        type="is-dark"
                        class="btn-new"
                        icon-left="check-circle-outline"
                        :loading="busyId === d.id"
                        :disabled="isBusy"
                        @click="onVerifyDomain(d)"
                      >
                        Verify DNS
                      </b-button>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <b-loading :is-full-page="false" :active="isBusy && domains.length === 0 && activeTab === 'domains'" />
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import CopyText from '../components/CopyText.vue';
import EmptyPlaceholder from '../components/EmptyPlaceholder.vue';

export default Vue.extend({
  components: {
    CopyText,
    EmptyPlaceholder,
  },

  data() {
    return {
      domains: [],
      busyId: null,
      expandedDomain: null,
      expandedSender: null,
      trackingHostDraft: {},
      domainQuery: '',
      senderQuery: '',
      fromEmail: '',
    };
  },

  methods: {
    setTab(tab) {
      const query = { ...this.$route.query, tab };
      this.$router.replace({ query }).catch(() => {});
    },

    onHeaderAdd() {
      if (this.activeTab === 'domains') {
        this.$router.push({ name: 'trackingDomainAdd' });
        return;
      }
      this.$router.push({ name: 'senderAdd' });
    },

    parseFrom(raw) {
      const s = (raw || '').trim();
      const m = s.match(/^(.*?)\s*<([^>]+)>$/);
      if (m) {
        return {
          name: m[1].replace(/^["']|["']$/g, '').trim(),
          email: m[2].trim(),
        };
      }
      if (s.indexOf('@') > -1) {
        return { name: '', email: s };
      }
      return { name: '', email: '' };
    },

    senderLabel(s) {
      if (s.name && s.email) return `${s.name} <${s.email}>`;
      return s.email || s.name || '';
    },

    emailDomain(email) {
      const parts = (email || '').split('@');
      return parts.length === 2 ? parts[1].toLowerCase() : '';
    },

    domainForEmail(email) {
      const host = this.emailDomain(email);
      if (!host) return null;
      return this.domains.find((d) => d.status === 'verified' && (
        d.domain === host || d.domain.endsWith(`.${host}`)
      )) || null;
    },

    dnsHostLabel(trackingHost, baseDomain) {
      const host = this.normalizeHost(trackingHost);
      const base = this.normalizeHost(baseDomain || trackingHost);
      if (!host) return '';
      if (host === base) return '';
      if (host.endsWith(`.${base}`)) {
        return host.slice(0, -(base.length + 1));
      }
      const parts = host.split('.').filter(Boolean);
      if (parts.length > 2) return parts[0];
      return host;
    },

    normalizeHost(raw) {
      return String(raw || '')
        .trim()
        .toLowerCase()
        .replace(/^https?:\/\//, '')
        .replace(/^www\./, '')
        .replace(/\/.*$/, '')
        .replace(/\.$/, '');
    },

    baseDomainFor(d) {
      return this.normalizeHost(d.baseDomain || d.domain);
    },

    trackingHostFor(d) {
      const draft = this.trackingHostDraft[d.id];
      if (draft !== undefined && draft !== null) {
        return this.normalizeHost(draft);
      }
      return this.normalizeHost(d.domain);
    },

    dnsHostLabelFor(d) {
      const host = this.trackingHostFor(d);
      const base = this.baseDomainFor(d);
      const label = this.dnsHostLabel(host, base);
      return label || '—';
    },

    isValidTrackingHost(d) {
      const host = this.trackingHostFor(d);
      const base = this.baseDomainFor(d);
      if (!host || !base || host === base) return false;
      return host.endsWith(`.${base}`) && this.isValidDomain(host);
    },

    isValidDomain(domain) {
      return /^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+$/.test(domain);
    },

    onTrackingHostInput(d) {
      this.trackingHostDraft = {
        ...this.trackingHostDraft,
        [d.id]: this.trackingHostDraft[d.id] || '',
      };
    },

    toggleSender(email) {
      this.expandedSender = this.expandedSender === email ? null : email;
    },

    openSenderForm(sender) {
      const query = sender
        ? { name: sender.name || '', email: sender.email || '' }
        : {};
      this.$router.push({ name: 'senderAdd', query });
    },

    loadSender() {
      this.$api.getSettings().then((data) => {
        this.fromEmail = (data && data['app.from_email']) || '';
      });
    },

    normalizeList(data) {
      if (Array.isArray(data)) {
        return data.map(this.normalizeItem);
      }
      if (data && Array.isArray(data.results)) {
        return data.results.map(this.normalizeItem);
      }
      return [];
    },

    normalizeItem(item) {
      const d = item || {};
      const status = ['pending', 'verified', 'failed'].indexOf(d.status) > -1
        ? d.status : 'pending';
      return {
        ...d,
        status,
        dnsRecordType: d.dnsRecordType || 'CNAME',
        dnsRecordName: d.dnsRecordName || d.domain || '',
        dnsRecordValue: d.dnsRecordValue || '',
        baseDomain: d.baseDomain || d.base_domain || d.domain || '',
        lastError: d.lastError || '',
      };
    },

    getDomains() {
      return this.$api.getTrackingDomains().then((data) => {
        this.domains = this.normalizeList(data);
      });
    },

    onDomainAction(d) {
      const opening = this.expandedDomain !== d.id;
      this.expandedDomain = opening ? d.id : null;
      if (opening) {
        const host = this.normalizeHost(d.domain);
        const base = this.baseDomainFor(d);
        this.trackingHostDraft = {
          ...this.trackingHostDraft,
          [d.id]: host !== base ? host : '',
        };
      }
    },

    onVerifyDomain(d) {
      if (!this.isValidTrackingHost(d)) {
        this.$utils.toast(
          `Enter a tracking subdomain such as emailtrack.${this.baseDomainFor(d)} before verifying.`,
          'is-warning',
          5000,
        );
        return;
      }

      this.busyId = d.id;
      const payload = { domain: this.trackingHostFor(d) };
      this.$api.verifyTrackingDomain(d.id, payload).then((data) => {
        this.busyId = null;
        const resp = data || {};
        const record = resp.id ? resp : resp.domain;
        const status = (record && record.status) || resp.status
          || (resp.verified === true ? 'verified' : null);

        if (record && record.id) {
          this.domains = this.domains.map((row) => (
            row.id === record.id ? this.normalizeItem({ ...row, ...record }) : row
          ));
        } else {
          this.getDomains();
        }

        if (status === 'verified' || resp.verified === true) {
          this.$utils.toast(resp.message || `${d.domain} is verified.`);
          return;
        }

        const err = resp.message || (record && record.lastError);
        this.$utils.toast(
          err || 'DNS record not found yet. DNS changes may take time to propagate.',
          'is-warning',
          5000,
        );
      }, () => {
        this.busyId = null;
      });
    },

    displayDomainName(d) {
      // Always show the registered base domain in the list.
      // The full tracking host (e.g. emailtrack.eguardian.in) is shown
      // only inside the expanded DNS configuration panel.
      return this.baseDomainFor(d) || this.normalizeHost(d.domain);
    },

    onDeleteDomain(d) {
      this.$utils.confirm(
        `Remove ${d.domain}? New campaigns will fall back to the platform domain. `
          + 'Tracking links in emails already sent still require this hostname and its DNS record.',
        () => {
          this.$api.deleteTrackingDomain(d.id).then(() => {
            if (this.expandedDomain === d.id) this.expandedDomain = null;
            this.getDomains();
            this.$utils.toast(this.$t('globals.messages.deleted', { name: d.domain }));
          });
        },
      );
    },
  },

  computed: {
    ...mapState(['loading']),

    activeTab() {
      return this.$route.query.tab === 'domains' ? 'domains' : 'senders';
    },

    isBusy() {
      return this.loading.trackingDomains;
    },

    canManage() {
      return this.$can('settings:manage');
    },

    senders() {
      const parsed = this.parseFrom(this.fromEmail);
      if (!parsed.email) return [];
      return [{
        name: parsed.name,
        email: parsed.email,
        domain: this.emailDomain(parsed.email),
      }];
    },

    filteredSenders() {
      const q = (this.senderQuery || '').trim().toLowerCase();
      if (!q) return this.senders;
      return this.senders.filter((s) => (
        `${s.name} ${s.email}`.toLowerCase().indexOf(q) > -1
      ));
    },

    filteredDomains() {
      const q = (this.domainQuery || '').trim().toLowerCase();
      if (!q) return this.domains;
      return this.domains.filter((d) => (d.domain || '').toLowerCase().indexOf(q) > -1);
    },
  },

  created() {
    this.$root.$on('page.refresh', this.getDomains);
  },

  destroyed() {
    this.$root.$off('page.refresh', this.getDomains);
  },

  mounted() {
    this.getDomains();
    this.loadSender();
  },
});
</script>
