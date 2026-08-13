<template>
  <form @submit.prevent="onSubmit">
    <div class="modal-card content" style="width: auto">
      <header class="modal-card-head">
        <p v-if="isEditing" class="has-text-grey-light is-size-7">
          {{ $t('globals.fields.id') }}: <copy-text :text="`${data.id}`" />
          {{ $t('globals.fields.uuid') }}: <copy-text :text="data.uuid" />
        </p>
        <h4 v-if="isEditing">
          {{ data.name }}
        </h4>
        <h4 v-else>
          {{ $t('forms.newForm') }}
        </h4>
      </header>
      <section expanded class="modal-card-body">
        <b-field :label="$t('globals.fields.name')" label-position="on-border">
          <b-input :maxlength="200" :ref="'focus'" v-model="form.name" name="name"
            :placeholder="$t('globals.fields.name')" required />
        </b-field>

        <list-selector v-model="selectedLists" :all="publicLists" :selected="selectedLists"
          :label="$t('forms.publicLists')" :placeholder="$t('forms.selectHelp')" />

        <b-field :label="$t('forms.redirectURL')" :message="$t('forms.redirectURLHelp')">
          <div class="content">
            <ul class="no">
              <li>
                <b-radio v-model="form.redirectUrl" native-value="">
                  {{ $t('globals.terms.none') }}
                </b-radio>
              </li>
              <li v-for="url in redirectURLs" :key="url">
                <b-radio v-model="form.redirectUrl" :native-value="url">
                  {{ url }}
                </b-radio>
              </li>
            </ul>
          </div>
        </b-field>
      </section>
      <footer class="modal-card-foot has-text-right">
        <b-button @click="$parent.close()">
          {{ $t('globals.buttons.close') }}
        </b-button>
        <b-button native-type="submit" type="is-primary" :loading="loading.forms" data-cy="btn-save">
          {{ $t('globals.buttons.save') }}
        </b-button>
      </footer>
    </div>
  </form>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import CopyText from '../components/CopyText.vue';
import ListSelector from '../components/ListSelector.vue';

export default Vue.extend({
  name: 'FormEdit',

  components: {
    CopyText,
    ListSelector,
  },

  props: {
    data: { type: Object, default: () => ({}) },
    isEditing: { type: Boolean, default: false },
  },

  data() {
    return {
      form: {
        name: '',
        redirectUrl: '',
      },
      selectedLists: [],
    };
  },

  methods: {
    onSubmit() {
      const data = {
        name: this.form.name,
        redirect_url: this.form.redirectUrl,
        list_ids: this.selectedLists.map((l) => l.id),
      };

      if (this.isEditing) {
        this.$api.updateForm({ id: this.data.id, ...data }).then((d) => {
          this.$emit('finished');
          this.$parent.close();
          this.$utils.toast(this.$t('globals.messages.updated', { name: d.name }));
        });
        return;
      }

      this.$api.createForm(data).then((d) => {
        this.$emit('finished');
        this.$parent.close();
        this.$utils.toast(this.$t('globals.messages.created', { name: d.name }));
      });
    },
  },

  computed: {
    ...mapState(['loading', 'lists', 'serverConfig']),

    publicLists() {
      if (!this.lists.results) {
        return [];
      }
      return this.lists.results.filter((l) => l.type === 'public');
    },

    redirectURLs() {
      const urls = this.serverConfig.public_subscription
        ? this.serverConfig.public_subscription.redirect_urls
        : [];
      return Array.isArray(urls) ? urls : [];
    },
  },

  mounted() {
    this.form = { ...this.form, ...this.$props.data };
    this.selectedLists = this.publicLists.filter(
      (l) => (this.$props.data.listIds || []).includes(l.id),
    );

    this.$nextTick(() => {
      this.$refs.focus.focus();
    });
  },
});
</script>
