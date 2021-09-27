<template>
  <div>
    <v-dialog
        v-model="playerDialog"
        @click:outside="closeDialog"
        max-width="400"
    >
      <v-card>
        <video-player class="video-player-box"
                      ref="videoPlayer"
                      :options="playerOptions"
                      :playsinline="true"
        >
        </video-player>
      </v-card>
    </v-dialog>

    <v-container>
      <h2 class="text-center pa-8">A-SOUL 抖音视频 All in one!</h2>
      <v-row
          justify="center"
          no-gutters
      >
        <v-text-field
            ref="searchInput"
            v-model="searchString"
            class="mx-2 mx-md-4 rounded-lg justify-center"
            placeholder="搜索..."
            autocomplete="off"
            dense
            hide-details
            solo
            style="max-width: 450px;"
            @input="searchChange"
        >
        </v-text-field>
      </v-row>
      <v-row>

      </v-row>
    </v-container>


    <v-container>
      <v-row dense>
        <v-col
            v-for="v in videos"
            :key="v.id"
            :xl="3"
            :lg="3"
            :md="3"
            :sm="12"
        >
          <v-card>
            <v-img
                :src="v.dynamic_cover_urls[0]"
                class="white--text align-end"
                gradient="to bottom, rgba(0,0,0,.1), rgba(0,0,0,.5)"
                height="220px"
            >
              <v-card-title v-text="v.description"></v-card-title>
            </v-img>

            <v-card-actions>
              <div>
                <v-icon color="red lighten-2">mdi-heart</v-icon>
                <span class="grey--text text--darken-1">{{ v.statistic.digg }}</span>

                <v-icon color="gray" class="hidden-md-and-up">mdi-share</v-icon>
                <span class="grey--text text--darken-1 hidden-md-and-up"> {{ v.statistic.share }}</span>

                <v-icon color="gray" class="hidden-md-and-up">mdi-comment</v-icon>
                <span class="grey--text text--darken-1 hidden-md-and-up"> {{ v.statistic.comment }}</span>
              </div>

              <v-spacer></v-spacer>
              <span class="grey--text text--darken-1">{{ new Date(v.created_at).toLocaleDateString() }}</span>
              <v-spacer></v-spacer>
              <v-btn icon @click="playVideo(v)">
                <v-icon>mdi-play</v-icon>
              </v-btn>
              <v-btn icon>
                <v-icon @click="downloadVideo(v)">mdi-download</v-icon>
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
        <v-col
            :xl="3"
            :lg="3"
            :md="3"
            :sm="12"
        >
          <v-skeleton-loader
              v-if="isLoading"
              class="mx-auto"
              type="card"
          ></v-skeleton-loader>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import axios from 'axios'
import 'video.js/dist/video-js.css'
import {videoPlayer} from 'vue-video-player'

export default {
  name: "Index",

  data() {
    return {
      isLoading: true,
      isEnd: false,
      currentPage: 1,
      playerDialog: false,
      playerOptions: {
        loop: true,
        fluid: true,
        autoplay: true,
        muted: false,
        playbackRates: [0.7, 1.0, 1.5, 2.0],
        sources: [{
          type: "video/mp4",
          src: ''
        }]
      },
      videos: [],

      searchString: ''
    }
  },

  mounted() {
    this.getVideos()
    this.onScroll()
  },

  methods: {
    getVideos() {
      return new Promise((resolve, reject) => {
        if (this.isEnd) {
          resolve()
          return
        }
        this.isLoading = true
        axios.get('https://asoul.cdn.n3ko.co/api/videos', {
          params: {
            page: this.currentPage,
            keyword: this.searchString,
          }
        }).then(res => {
          if (res.data.data.length === 0) {
            this.isEnd = true
          }

          this.videos = this.videos.concat(res.data.data)
          this.isLoading = false
          resolve()
        }).catch(err => {
          this.isLoading = false
          reject(err)
        })
      })
    },

    playVideo(v) {
      this.playerDialog = true
      this.playerOptions.sources[0].src = v.video_urls[0]
      this.playerOptions.sources[0].poster = v.dynamic_cover_urls[0]
      if (this.$refs.videoPlayer !== undefined) {
        this.$refs.videoPlayer.player.play()
      }
    },

    closeDialog() {
      this.$refs.videoPlayer.player.pause()
    },

    downloadVideo(v) {
      window.open(v.video_urls[0])
    },

    onScroll() {
      let isLoading = false
      window.onscroll = () => {
        let bottomOfWindow = document.documentElement.offsetHeight - document.documentElement.scrollTop - window.innerHeight <= 200
        if (bottomOfWindow && isLoading === false) {
          isLoading = true
          this.currentPage++
          this.getVideos().then(() => {
            isLoading = false
          }).catch(() => {
            isLoading = false
          })
        }
      }
    },

    searchChange() {
      this.isEnd = false
      this.videos = []
      this.currentPage = 1
      this.getVideos()
    },

    resetSearch() {
      this.searchString = ''
    },
  },

  components: {
    videoPlayer
  }
}
</script>

<style scoped>

</style>