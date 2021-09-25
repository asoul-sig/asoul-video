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
      <v-row dense>
        <v-col
            v-for="v in videos"
            :key="v.id"
            :cols="3"
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
      videos: []
    }
  },

  mounted() {
    this.getVideos(this.currentPage)
    this.onScroll()
  },

  methods: {
    getVideos(page) {
      return new Promise((resolve, reject) => {
        axios.get('/api/videos?page=' + page).then(res => {
          this.videos = this.videos.concat(res.data.data)
          resolve()
        }).catch(err => {
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
          this.getVideos(++this.currentPage).then(() => {
            isLoading = false
          }).catch(() => {
            isLoading = false
          })
        }
      }
    }
  },

  components: {
    videoPlayer
  }
}
</script>

<style scoped>

</style>