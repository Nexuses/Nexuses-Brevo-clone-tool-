<template>
  <section class="crm-page tasks-page tasks-brevo">
    <crm-subnav />

    <header class="tasks-brevo__header">
      <div class="tasks-brevo__header-main">
        <h1 class="tasks-brevo__title">Tasks</h1>
        <p class="tasks-brevo__lead">
          Create and manage tasks for follow-ups, calls, and to-dos across your CRM.
        </p>
      </div>
      <button type="button" class="tasks-brevo__create" @click="openCreate">
        <span class="tasks-brevo__create-plus" aria-hidden="true">+</span>
        Create a task
      </button>
    </header>

    <div class="crm-toolbar">
      <b-field>
        <b-input
          v-model="query"
          expanded
          placeholder="Search a task name"
          icon="magnify"
        />
      </b-field>
    </div>

    <div class="crm-card-table">
      <table class="table is-fullwidth is-hoverable">
        <thead>
          <tr>
            <th>Task</th>
            <th>Type</th>
            <th>Due date</th>
            <th>Priority</th>
            <th>Status</th>
            <th class="has-text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in filteredTasks" :key="t.id">
            <td>
              <a href="#" class="crm-link" @click.prevent="openEdit(t)">{{ t.name }}</a>
            </td>
            <td class="is-capitalized">{{ typeLabel(t.type) }}</td>
            <td>{{ formatDue(t) }}</td>
            <td>
              <span v-if="t.highPriority" class="task-priority is-high">High</span>
              <span v-else class="has-text-grey">Normal</span>
            </td>
            <td class="is-capitalized">{{ t.status || 'open' }}</td>
            <td class="has-text-right">
              <b-dropdown
                class="bv-action-menu"
                position="is-bottom-left"
                :mobile-modal="false"
              >
                <template #trigger>
                  <button type="button" class="campaign-actions-trigger" aria-label="Actions">
                    <span class="campaign-kebab" aria-hidden="true">
                      <span /><span /><span />
                    </span>
                  </button>
                </template>
                <b-dropdown-item @click="openEdit(t)">
                  <b-icon icon="pencil-outline" size="is-small" />
                  Edit
                </b-dropdown-item>
                <b-dropdown-item v-if="t.status !== 'done'" @click="completeTask(t)">
                  <b-icon icon="check-circle-outline" size="is-small" />
                  Complete
                </b-dropdown-item>
                <hr class="bv-action-menu__sep" />
                <b-dropdown-item class="bv-action-menu__danger" @click="deleteTask(t)">
                  <b-icon icon="trash-can-outline" size="is-small" />
                  Delete
                </b-dropdown-item>
              </b-dropdown>
            </td>
          </tr>
          <tr v-if="filteredTasks.length === 0">
            <td colspan="6" class="has-text-centered has-text-grey py-5">
              No tasks yet. Create one to get started.
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <b-modal :active.sync="isFormVisible" :width="560" scroll="keep" class="task-form-modal">
      <task-form
        v-if="isFormVisible"
        :data="curItem"
        @close="isFormVisible = false"
        @finished="onFinished"
      />
    </b-modal>
  </section>
</template>

<script>
import dayjs from 'dayjs';
import CrmSubnav from '../components/CrmSubnav.vue';
import TaskForm, { loadTasks, saveTasks } from '../components/TaskForm.vue';

export default {
  name: 'Tasks',
  components: { CrmSubnav, TaskForm },

  data() {
    return {
      query: '',
      tasks: [],
      isFormVisible: false,
      curItem: null,
    };
  },

  computed: {
    filteredTasks() {
      const q = this.query.trim().toLowerCase();
      if (!q) return this.tasks;
      return this.tasks.filter((t) => (t.name || '').toLowerCase().includes(q));
    },
  },

  methods: {
    refresh() {
      this.tasks = loadTasks();
    },

    typeLabel(type) {
      const map = {
        todo: 'To do', call: 'Call', email: 'Email', meeting: 'Meeting',
      };
      return map[type] || type || 'To do';
    },

    formatDue(t) {
      if (!t.dueDate) return '-';
      const date = dayjs(t.dueDate).format('MMM D, YYYY');
      return t.time ? `${date} - ${t.time}` : date;
    },

    openCreate() {
      this.curItem = null;
      this.isFormVisible = true;
    },

    openEdit(t) {
      this.curItem = { ...t };
      this.isFormVisible = true;
    },

    onFinished() {
      this.isFormVisible = false;
      this.refresh();
    },

    completeTask(t) {
      const next = loadTasks().map((x) => (
        x.id === t.id ? { ...x, status: 'done', updatedAt: new Date().toISOString() } : x
      ));
      saveTasks(next);
      this.refresh();
      this.$utils.toast('Task completed');
    },

    deleteTask(t) {
      this.$utils.confirm(`Delete task "${t.name}"?`, () => {
        saveTasks(loadTasks().filter((x) => x.id !== t.id));
        this.refresh();
        this.$utils.toast('Task deleted');
      });
    },
  },

  mounted() {
    this.refresh();
    if (this.$route.query.create === '1') {
      this.openCreate();
    }
  },
};
</script>
