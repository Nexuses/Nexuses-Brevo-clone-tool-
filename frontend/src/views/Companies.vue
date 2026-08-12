<template>
  <section class="crm-page companies-page companies-brevo">
    <crm-subnav />

    <header class="companies-brevo__header">
      <div class="companies-brevo__header-main">
        <h1 class="companies-brevo__title">Companies</h1>
        <p class="companies-brevo__lead">
          Store company records linked to your contacts for account-based outreach.
        </p>
      </div>
      <div class="companies-brevo__actions">
        <button type="button" class="bv-btn-outline" data-cy="btn-import-companies" @click="importHint">
          Import companies
        </button>
        <button type="button" class="companies-brevo__create" data-cy="btn-new-company" @click="createCompany">
          <span class="companies-brevo__create-plus" aria-hidden="true">+</span>
          Create company
        </button>
      </div>
    </header>

    <div class="crm-view-tabs">
      <button type="button" class="crm-view-tabs__tab is-active">All companies</button>
    </div>

    <div class="crm-filter-bar">
      <b-dropdown>
        <template #trigger>
          <button type="button" class="crm-filter-bar__btn">
            Add filter <b-icon icon="arrow-down" size="is-small" />
          </button>
        </template>
        <b-dropdown-item disabled>Coming soon</b-dropdown-item>
      </b-dropdown>
    </div>

    <div class="crm-meta-row">
      <span class="crm-pill">{{ filteredCompanies.length }} companies</span>
      <div class="crm-meta-row__right">
        <b-field>
          <b-input
            v-model="query"
            expanded
            placeholder="Company name, domain"
            icon="magnify"
            data-cy="company-search"
          />
        </b-field>
      </div>
    </div>

    <div class="crm-card-table">
      <table class="table is-fullwidth is-hoverable">
        <thead>
          <tr>
            <th style="width: 40px">
              <b-checkbox v-model="allCheckedProxy" />
            </th>
            <th>Company name</th>
            <th>Domain</th>
            <th>Owner</th>
            <th>Phone number</th>
            <th class="has-text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in filteredCompanies" :key="c.id">
            <td>
              <b-checkbox v-model="checked" :native-value="c.id" />
            </td>
            <td>
              <a href="#" class="crm-company-name" @click.prevent="editCompany(c)">{{ c.name }}</a>
            </td>
            <td>
              <a v-if="c.domain" :href="`https://${c.domain}`" target="_blank" rel="noopener noreferrer" class="crm-link">
                {{ c.domain }}
              </a>
              <span v-else class="has-text-grey">--</span>
            </td>
            <td>{{ c.owner || '--' }}</td>
            <td>{{ c.phone || '--' }}</td>
            <td class="has-text-right">
              <b-dropdown position="is-bottom-left">
                <template #trigger>
                  <button type="button" class="crm-kebab" aria-label="Actions">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                      <circle cx="8" cy="3.5" r="1.4" fill="currentColor" />
                      <circle cx="8" cy="8" r="1.4" fill="currentColor" />
                      <circle cx="8" cy="12.5" r="1.4" fill="currentColor" />
                    </svg>
                  </button>
                </template>
                <b-dropdown-item @click="editCompany(c)">Edit</b-dropdown-item>
                <b-dropdown-item class="has-text-danger" @click="deleteCompany(c)">
                  Delete
                </b-dropdown-item>
              </b-dropdown>
            </td>
          </tr>
          <tr v-if="filteredCompanies.length === 0">
            <td colspan="6" class="has-text-centered has-text-grey py-5">
              No companies yet. Create one to get started.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script>
import CrmSubnav from '../components/CrmSubnav.vue';

const STORAGE_KEY = 'nexuses.crm.companies';

const DEFAULT_COMPANIES = [
  {
    id: 1,
    name: 'Acme Inc',
    domain: 'acme.example',
    owner: '',
    phone: '',
  },
];

export default {
  name: 'Companies',
  components: { CrmSubnav },

  data() {
    return {
      query: '',
      companies: [],
      checked: [],
    };
  },

  computed: {
    filteredCompanies() {
      const q = this.query.trim().toLowerCase();
      if (!q) return this.companies;
      return this.companies.filter((c) => (
        c.name.toLowerCase().includes(q)
        || (c.domain || '').toLowerCase().includes(q)
        || (c.owner || '').toLowerCase().includes(q)
      ));
    },

    allChecked() {
      return this.filteredCompanies.length > 0
        && this.filteredCompanies.every((c) => this.checked.includes(c.id));
    },

    allCheckedProxy: {
      get() {
        return this.allChecked;
      },
      set(val) {
        this.toggleAll(val);
      },
    },
  },

  methods: {
    load() {
      try {
        const raw = localStorage.getItem(STORAGE_KEY);
        this.companies = raw ? JSON.parse(raw) : [...DEFAULT_COMPANIES];
      } catch (e) {
        this.companies = [...DEFAULT_COMPANIES];
      }
    },

    save() {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(this.companies));
    },

    toggleAll(val) {
      if (val) {
        this.checked = this.filteredCompanies.map((c) => c.id);
      } else {
        this.checked = [];
      }
    },

    importHint() {
      this.$utils.toast('Company import coming soon');
    },

    createCompany() {
      this.$utils.prompt('Create company', {
        placeholder: 'Company name',
        value: '',
      }, (name) => {
        const trimmed = (name || '').trim();
        if (!trimmed) return;
        const nextId = this.companies.reduce((max, c) => Math.max(max, c.id), 0) + 1;
        this.companies.unshift({
          id: nextId,
          name: trimmed,
          domain: '',
          owner: '',
          phone: '',
        });
        this.save();
        this.$utils.toast('Company created');
      });
    },

    editCompany(c) {
      this.$utils.prompt('Edit company', {
        placeholder: 'Company name',
        value: c.name,
      }, (name) => {
        const trimmed = (name || '').trim();
        if (!trimmed) return;
        this.companies = this.companies.map((item) => (
          item.id === c.id ? { ...item, name: trimmed } : item
        ));
        this.save();
        this.$utils.toast('Company updated');
      });
    },

    deleteCompany(c) {
      this.$utils.confirm(`Delete company "${c.name}"?`, () => {
        this.companies = this.companies.filter((x) => x.id !== c.id);
        this.checked = this.checked.filter((id) => id !== c.id);
        this.save();
        this.$utils.toast('Company deleted');
      });
    },
  },

  mounted() {
    this.load();
  },
};
</script>
