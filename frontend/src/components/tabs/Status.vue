<template>
  <v-container>
    <v-alert
      v-if="profile.key === ''"
      dense
      outlined
      type="error"
    >No profile set, cannot query API</v-alert>
    <v-row>
      <v-col>
        Auto refresh is currently 
        <span
          color="error"
          v-if="layout && layout.auto"
        >ON</span>
        <span v-else
        >OFF</span>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        Last refresh was {{ updateDuration }} ago
      </v-col>
    </v-row>
    <v-row justify="space-around">
      <v-col cols="auto">
        <v-btn
          @click="update()"
          :loading="loading"
          color="primary"
        >Refresh Now</v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapState, mapActions } from 'vuex';

export default {
  data: () => ({
    loading: false,
  }),

  computed: {
    ...mapState('profile', ['profile']),
    ...mapState('layout', ['layout']),
    ...mapState('game', ['handle', 'timer']),

    updateDuration() {
      let seconds = this.timer;
      const hours = Math.floor(seconds / 3600);
      seconds -= hours * 3600;
      const minutes = Math.floor(seconds / 60);
      seconds -= minutes * 60;

      const hourChunk = hours > 0 ? `${hours}h ` : '';
      const minuteChunk = minutes > 0 || hours > 0 ? `${minutes}m ` : '';
      return `${hourChunk}${minuteChunk}${seconds}s`;
    },
  },

  methods: {
    ...mapActions('game', [
      'getGameInfo',
      'setRefreshTicker', 'clearRefreshTicker',
      'incrementTimer', 'resetTimer',
    ]),
    tick() {
      this.incrementTimer();
      if (!this.layout || !this.layout.auto) return;
    },
    async update() {
      this.loading = true;
      await this.getGameInfo(this.profile);
      this.loading = false;
      this.resetTimer();
    },
  },

  mounted() {
    if (this.handle !== null) this.clearRefreshTicker();
    this.setRefreshTicker(setInterval(this.tick, 1000));
  },
};
</script>
