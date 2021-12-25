<template>
  <v-app>
    <v-app-bar
        app
        color="primary"
    >
      <div class="d-flex align-center white--text">
        A-SOUL Video
      </div>
      <v-spacer></v-spacer>
      <v-btn
          href="https://github.red/asoul-video-trick/"
          target="_blank"
          text
          class="white--text"
      >
        关于
      </v-btn>
      <v-menu offset-y open-on-hover>
        <template v-slot:activator="{ on, attrs }">
          <v-btn
              dark
              text
              v-bind="attrs"
              v-on="on"
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
      <v-btn
          href="https://github.com/asoul-video"
          target="_blank"
          text
          class="white--text"
      >
        <span class="mr-2">GitHub</span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
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
