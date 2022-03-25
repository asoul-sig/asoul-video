<template>
  <v-app>
    <v-app-bar
        app
        color="primary"
    >
      <v-spacer></v-spacer>
      <v-spacer></v-spacer>
      <v-btn
          href="https://github.red/asoul-video-trick/"
          target="_blank"
          text
          class="white--text ml-2 mr-2"
      >
        关于
      </v-btn>
      <v-btn
          href="https://github.com/asoul-sig"
          target="_blank"
          text
          class="white--text ml-2 mr-2"
      >
        <span>开源</span>
      </v-btn>
      <div class="white--text ml-2 mr-2">
        <b>A-SOUL Video</b>
      </div>
      <v-btn
          href="https://github.com/asoul-video/acao"
          target="_blank"
          text
          class="white--text ml-2 mr-2"
      >
        阿草
      </v-btn>
      <v-menu offset-y open-on-hover>
        <template v-slot:activator="{ on, attrs }">
          <v-btn
              dark
              text
              v-bind="attrs"
              v-on="on"
              class="white--text ml-2 mr-2"
          >
            友链
          </v-btn>
        </template>
        <v-list>
          <v-list-item-group>
            <v-list-item v-for="(link, index) in friendLink" v-bind:key="index">
              <v-list-item-action @click="open(link.url)">
                {{ link.name }}
              </v-list-item-action>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-menu>
      <v-spacer></v-spacer>
      <v-spacer></v-spacer>
    </v-app-bar>
    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>


<script>
export default {
  name: 'App',

  data() {
    return {
      friendLink: [{"name": "枝网查重", "url": "https://asoulcnki.asia/"}]
    }
  },

  methods: {
    open(url) {
      window.open(url, '_blank')
    }
  }
};
</script>
