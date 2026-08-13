<template>
  <div class="modal-card content" style="width: auto">
    <header class="modal-card-head">
      <h4>{{ data.name }}</h4>
    </header>
    <section expanded class="modal-card-body">
      <div class="columns is-vcentered form-url">
        <div class="column">
          <a :href="pageURL" target="_blank" rel="noopener noreferer" data-cy="url">{{ pageURL }}</a>
        </div>
        <div class="column is-narrow">
          <b-button icon-left="file-multiple-outline" @click="copyLink" data-cy="btn-copy-url">
            {{ $t('forms.copyLink') }}
          </b-button>
        </div>
      </div>

      <b-tabs>
        <b-tab-item :label="$t('forms.formHTML')">
          <p>
            {{ $t('forms.formHTMLHelp') }}
          </p>
          <code-editor lang="html" v-model="html" disabled />
        </b-tab-item>

        <b-tab-item :label="$t('forms.formIframe')">
          <p>
            {{ $t('forms.formIframeHelp') }}
          </p>
          <code-editor lang="html" v-model="iframe" disabled />
        </b-tab-item>
      </b-tabs>
    </section>
    <footer class="modal-card-foot has-text-right">
      <b-button @click="$parent.close()">
        {{ $t('globals.buttons.close') }}
      </b-button>
    </footer>
  </div>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import CodeEditor from '../components/CodeEditor.vue';

export default Vue.extend({
  name: 'FormEmbed',

  components: {
    'code-editor': CodeEditor,
  },

  props: {
    data: { type: Object, default: () => ({}) },
  },

  methods: {
    escapeAttr(value) {
      return String(value)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
    },

    copyLink() {
      const input = document.createElement('input');
      input.setAttribute('type', 'text');
      input.style.opacity = '0';
      input.value = this.pageURL;
      document.body.appendChild(input);
      input.select();
      document.execCommand('copy');
      document.body.removeChild(input);

      this.$utils.toast(this.$t('globals.messages.copied'));
    },
  },

  computed: {
    ...mapState(['lists', 'serverConfig']),

    formLists() {
      if (!this.lists.results) {
        return [];
      }
      const ids = this.data.listIds || [];
      return this.lists.results.filter((l) => ids.includes(l.id));
    },

    pageURL() {
      return `${this.serverConfig.root_url}/subscription/form?form=${this.data.uuid}`;
    },

    html() {
      let h = `<form method="post" action="${this.serverConfig.root_url}/subscription/form" class="listmonk-form">\n`
        + '  <div>\n'
        + `    <h3>${this.$t('public.sub')}</h3>\n`
        + '    <input type="hidden" name="nonce" />\n';

      if (this.data.redirectUrl) {
        h += `    <input type="hidden" name="next" value="${this.escapeAttr(this.data.redirectUrl)}" />\n`;
      }

      h += '\n'
        + `    <p><input type="email" name="email" required placeholder="${this.$t('subscribers.email')}" /></p>\n`
        + `    <p><input type="text" name="name" placeholder="${this.$t('public.subName')}" /></p>\n\n`;

      this.formLists.forEach((l) => {
        h += '    <p>\n'
          + `      <input id="${l.uuid.substr(0, 5)}" type="checkbox" name="l" checked value="${l.uuid}" />\n`
          + `      <label for="${l.uuid.substr(0, 5)}">${l.name}</label>\n`;

        if (l.description) {
          h += '      <br />\n'
            + `      <span>${l.description}</span>\n`;
        }

        h += '    </p>\n';
      });

      if (this.serverConfig.public_subscription.captcha_enabled) {
        if (this.serverConfig.public_subscription.captcha_provider === 'altcha') {
          h += '\n'
            + `    <altcha-widget challengeurl="${this.serverConfig.root_url}/api/public/captcha/altcha"></altcha-widget>\n`
            + `    <${'script'} type="module" src="${this.serverConfig.root_url}/public/static/altcha.umd.js" async defer></${'script'}>\n`;
        } else if (this.serverConfig.public_subscription.captcha_provider === 'hcaptcha') {
          h += '\n'
            + `    <div class="h-captcha" data-sitekey="${this.serverConfig.public_subscription.captcha_key}"></div>\n`
            + `    <${'script'} src="https://js.hcaptcha.com/1/api.js" async defer></${'script'}>\n`;
        }
      }

      h += '\n'
        + `    <input type="submit" value="${this.$t('public.sub')} " />\n`
        + '  </div>\n'
        + '</form>';

      return h;
    },

    iframe() {
      return `<iframe src="${this.escapeAttr(this.pageURL)}&simple=1" width="100%" height="400" frameborder="0" scrolling="auto"></iframe>`;
    },
  },
});
</script>
