<template>
  <section class="lists">
    <crm-subnav />

    <header class="columns page-header lists-brevo__header">
      <div class="column is-8">
        <h1 class="title is-4 mb-2">Lists</h1>
        <p class="crm-page__lead">
          This is where you organize your lists. Create, modify, and manage custom lists for targeted
          interactions, and keep them in folders for easy navigation.
        </p>
        <div class="is-size-7 mt-2">
          <router-link v-if="queryParams.status !== 'archived'" :to="{ name: 'lists', query: { status: 'archived' } }">
            {{ $t('globals.buttons.view') }} {{ $t('lists.archived').toLowerCase() }} &rarr;
          </router-link>
          <router-link v-else :to="{ name: 'lists' }">
            {{ $t('globals.buttons.view') }} {{ $t('menu.allLists').toLowerCase() }} &rarr;
          </router-link>
        </div>
      </div>
      <div class="column has-text-right">
        <b-button
          v-if="$can('lists:manage_all')"
          type="is-dark"
          icon-left="plus"
          class="btn-new"
          @click="showNewForm"
          data-cy="btn-new"
        >
          Create a list
        </b-button>
      </div>
    </header>

    <div class="crm-toolbar">
      <div class="crm-toolbar__folder">
        All folders ({{ lists.total || 0 }} lists)
      </div>
      <b-field>
        <b-input
          v-model="queryParams.query"
          name="query"
          expanded
          icon="magnify"
          placeholder="Search a list name or ID"
          ref="query"
          data-cy="query"
          @keyup.native.enter="getLists"
        />
      </b-field>
    </div>

    <b-table :data="lists.results" :loading="loading.listsFull" @check-all="onTableCheck" @check="onTableCheck"
      :checked-rows.sync="bulk.checked" hoverable default-sort="createdAt" paginated backend-pagination
      pagination-position="both" @page-change="onPageChange" :current-page="queryParams.page" :per-page="lists.perPage"
      :total="lists.total" checkable backend-sorting @sort="onSort" class="lists-brevo__table">
      <template #top-left>
        <div class="actions" v-if="bulk.checked.length > 0">
          <a class="a" href="#" @click.prevent="deleteLists" data-cy="btn-delete-lists">
            <b-icon icon="trash-can-outline" size="is-small" /> {{ $t('globals.buttons.delete') }}
          </a>
          <span class="a">
            {{ $tc('globals.messages.numSelected', numSelectedLists, { num: numSelectedLists }) }}
          </span>
        </div>
      </template>

      <b-table-column v-slot="props" field="name" label="Lists" header-class="cy-name" sortable
        width="28%" :td-attrs="$utils.tdID">
        <router-link class="crm-link" :to="`/contacts/lists/${props.row.id}`">
          {{ props.row.name }}
        </router-link>
      </b-table-column>

      <b-table-column v-slot="props" field="id" label="ID" header-class="cy-id" sortable width="8%">
        #{{ props.row.id }}
      </b-table-column>

      <b-table-column v-slot="props" field="type" label="Folder" header-class="cy-type" sortable width="12%">
        {{ $t(`lists.types.${props.row.type}`) }}
      </b-table-column>

      <b-table-column v-slot="props" field="subscriber_count" label="Contacts"
        header-class="cy-subscribers" numeric sortable centered>
        <template v-if="$can('subscribers:get_all', 'subscribers:get')">
          <router-link :to="`/contacts/lists/${props.row.id}`">
            {{ $utils.formatNumber(props.row.subscriberCount) }}
          </router-link>
        </template>
        <template v-else>
          {{ $utils.formatNumber(props.row.subscriberCount) }}
        </template>
      </b-table-column>

      <b-table-column v-slot="props" field="created_at" label="Creation date"
        header-class="cy-created_at" sortable>
        {{ $utils.niceDate(props.row.createdAt, true) }}
      </b-table-column>

      <b-table-column v-slot="props" label="Actions" cell-class="actions" align="right">
        <b-dropdown position="is-bottom-left" class="campaign-actions-menu">
          <template #trigger>
            <button type="button" class="campaign-actions-trigger" aria-label="Actions">
              <span class="campaign-kebab" aria-hidden="true"><span /><span /><span /></span>
            </button>
          </template>
          <div class="campaign-actions-panel">
            <router-link
              v-if="$can('subscribers:get_all', 'subscribers:get')"
              :to="`/contacts/lists/${props.row.id}`"
              class="campaign-action"
              aria-label="View contacts"
            >
              <b-tooltip label="View contacts" type="is-dark" position="is-left">
                <b-icon icon="account-multiple" />
              </b-tooltip>
            </router-link>
            <router-link
              v-if="$can('subscribers:manage')"
              :to="{ path: `/contacts/lists/${props.row.id}`, query: { add: '1' } }"
              class="campaign-action"
              aria-label="Add contact"
            >
              <b-tooltip label="Add a contact" type="is-dark" position="is-left">
                <b-icon icon="plus" />
              </b-tooltip>
            </router-link>
            <router-link
              v-if="$can('subscribers:import')"
              :to="{ name: 'import', query: { list_id: props.row.id } }"
              class="campaign-action"
              aria-label="Import contacts"
            >
              <b-tooltip label="Import contacts" type="is-dark" position="is-left">
                <b-icon icon="file-upload-outline" />
              </b-tooltip>
            </router-link>
            <a
              v-if="$can('lists:manage') || $canList(props.row.id, 'list:manage')"
              href="#"
              class="campaign-action"
              aria-label="Edit"
              @click.prevent="showEditForm(props.row)"
            >
              <b-tooltip :label="$t('globals.buttons.edit')" type="is-dark" position="is-left">
                <b-icon icon="pencil-outline" />
              </b-tooltip>
            </a>
            <router-link
              v-if="$can('campaigns:manage')"
              :to="`/campaigns/new?list_id=${props.row.id}`"
              class="campaign-action"
              aria-label="Campaign"
            >
              <b-tooltip :label="$t('lists.sendCampaign')" type="is-dark" position="is-left">
                <b-icon icon="rocket-launch-outline" />
              </b-tooltip>
            </router-link>
            <a
              v-if="$can('lists:manage') || $canList(props.row.id, 'list:manage')"
              href="#"
              class="campaign-action"
              aria-label="Delete"
              @click.prevent="deleteList(props.row)"
            >
              <b-tooltip :label="$t('globals.buttons.delete')" type="is-dark" position="is-left">
                <b-icon icon="trash-can-outline" />
              </b-tooltip>
            </a>
          </div>
        </b-dropdown>
      </b-table-column>

      <template #empty v-if="!loading.listsFull">
        <empty-placeholder />
      </template>
    </b-table>

    <!-- Add / edit form modal -->
    <b-modal scroll="keep" :aria-modal="true" :active.sync="isFormVisible" :width="600" @close="onFormClose">
      <list-form :data="curItem" :is-editing="isEditing" @finished="formFinished" />
    </b-modal>

    <p v-if="settings['app.cache_slow_queries']" class="has-text-grey">
      *{{ $t('globals.messages.slowQueriesCached') }}
      <a href="https://listmonk.app/docs/maintenance/performance/" target="_blank" rel="noopener noreferer"
        class="has-text-grey">
        <b-icon icon="link-variant" /> {{ $t('globals.buttons.learnMore') }}
      </a>
    </p>
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import EmptyPlaceholder from '../components/EmptyPlaceholder.vue';
import CrmSubnav from '../components/CrmSubnav.vue';
import ListForm from './ListForm.vue';

export default Vue.extend({
  components: {
    ListForm,
    EmptyPlaceholder,
    CrmSubnav,
  },

  data() {
    return {
      // Current list item being edited.
      curItem: null,
      isEditing: false,
      isFormVisible: false,
      lists: [],
      queryParams: {
        page: 1,
        query: '',
        orderBy: 'id',
        order: 'asc',
        status: this.$route.query.status || 'active',
      },

      // Table bulk row selection states.
      bulk: {
        checked: [],
        all: false,
      },
    };
  },

  methods: {
    onPageChange(p) {
      this.queryParams.page = p;
      this.getLists();
    },

    onSort(field, direction) {
      this.queryParams.orderBy = field;
      this.queryParams.order = direction;
      this.getLists();
    },

    // Show the edit list form.
    showEditForm(list) {
      this.curItem = list;
      this.isFormVisible = true;
      this.isEditing = true;
    },

    // Show the new list form.
    showNewForm() {
      this.curItem = {};
      this.isFormVisible = true;
      this.isEditing = false;
    },

    formFinished() {
      this.getLists();
    },

    onFormClose() {
      if (this.$route.params.id) {
        this.$router.push({ name: 'lists' });
      }
    },

    filterStatuses(list) {
      const out = { ...list.subscriberStatuses };
      if (list.optin === 'single') {
        delete out.unconfirmed;
        delete out.confirmed;
      }
      return out;
    },

    getLists() {
      this.$api.queryLists({
        page: this.queryParams.page,
        query: this.queryParams.query.replace(/[^\p{L}\p{N}\s]/gu, ' '),
        order_by: this.queryParams.orderBy,
        order: this.queryParams.order,
        status: this.queryParams.status,
      }).then((resp) => {
        this.lists = resp;
      });

      // Also fetch the minimal lists for the global store that appears
      // in dropdown menus on other pages like import and campaigns.
      this.$api.getLists({ minimal: true, per_page: 'all', status: 'active' });
    },

    deleteList(list) {
      this.$utils.confirm(
        this.$t('lists.confirmDelete'),
        () => {
          this.$api.deleteList(list.id).then(() => {
            this.getLists();

            this.$utils.toast(this.$t('globals.messages.deleted', { name: list.name }));
          });
        },
      );
    },

    // Mark all lists in the query as selected.
    onSelectAll() {
      this.bulk.all = true;
    },

    onTableCheck() {
      // Disable bulk.all selection if there are no rows checked in the table.
      if (this.bulk.checked.length !== this.lists.total) {
        this.bulk.all = false;
      }
    },

    deleteLists() {
      const name = this.$tc('globals.terms.list', this.numSelectedCampaigns);

      const fn = () => {
        const params = {};
        if (!this.bulk.all && this.bulk.checked.length > 0) {
          // If 'all' is not selected, delete lists by IDs.
          params.id = this.bulk.checked.map((l) => l.id);
        } else {
          // 'All' is selected, delete by query.
          params.query = this.queryParams.query.replace(/[^\p{L}\p{N}\s]/gu, ' ');
          params.all = this.bulk.all;
        }

        this.$api.deleteLists(params)
          .then(() => {
            this.getLists();
            this.$utils.toast(this.$tc(
              'globals.messages.deletedCount',
              this.numSelectedLists,
              { num: this.numSelectedLists, name },
            ));
          });
      };

      this.$utils.confirm(this.$tc(
        'globals.messages.confirmDelete',
        this.numSelectedLists,
        { num: this.numSelectedLists, name: name.toLowerCase() },
      ), fn);
    },

    createOptinCampaign(list) {
      const data = {
        name: this.$t('lists.optinTo', { name: list.name }),
        subject: this.$t('lists.confirmSub', { name: list.name }),
        lists: [list.id],
        from_email: this.settings['app.from_email'],
        content_type: 'richtext',
        messenger: 'email',
        type: 'optin',
      };

      this.$api.createCampaign(data).then((d) => {
        this.$router.push({ name: 'campaign', hash: '#content', params: { id: d.id } });
      });
      return false;
    },
  },

  computed: {
    ...mapState(['loading', 'settings']),

    numSelectedLists() {
      return this.bulk.all ? this.lists.total : this.bulk.checked.length;
    },
  },

  created() {
    this.$root.$on('page.refresh', this.getLists);
  },

  destroyed() {
    this.$root.$off('page.refresh', this.getLists);
  },

  mounted() {
    if (this.$route.params.id) {
      this.$api.getList(parseInt(this.$route.params.id, 10)).then((data) => {
        this.showEditForm(data);
      });
    } else {
      this.getLists();
    }
  },
});
</script>
