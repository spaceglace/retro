export default {
  namespaced: true,

  state: {
    game: null,
    gameMode: 'DateEarned',
    handle: null,
    timer: 0,
    order: [],
  },

  mutations: {
    SET_GAME_INFO(state, game) {
      state.game = game;
    },
    SET_HANDLE(state, handle) {
      state.handle = handle;
    },
    RESET_HANDLE(state) {
      clearInterval(state.handle);
      state.handle = null;
    },
    INCREMENT_TIMER(state) {
      state.timer += 1;
    },
    RESET_TIMER(state) {
      state.timer = 0;
    },
    SET_ORDER(state, order) {
      state.order = order;
    },
    SET_MODE(state, hardcore) {
      state.gameMode = hardcore ? 'DateEarnedHardcore' : 'DateEarned';
    },
  },

  actions: {
    async getGameInfo({ commit }, { name, key, dumb }) {
      try {
        const game = await window.backend.Retro.GetLastGame(name, key);
        const result = await window.backend.Retro.GetGameInformation(name, key, game);
        if (dumb) {
          for (let i in result.achievements) {
            if (!result.achievements[i].DateEarnedHardcore) {
              result.achievements[i].DateEarnedHardcore = "2020-04-28 23:42:48";
            }
          }
        }
        commit('SET_GAME_INFO', result);
      } catch (e) {
        console.error('Failed to get game info:', e);
      }
    },

    async setRefreshTicker({ commit }, handle) {
      commit('SET_HANDLE', handle);
    },
    async clearRefreshTicker({ commit }) {
      commit('RESET_HANDLE');
    },

    async resetTimer({ commit }) {
      commit('RESET_TIMER');
    },
    async incrementTimer({ commit }) {
      commit('INCREMENT_TIMER');
    },

    async setGameOrder({ commit }, { game, order }) {
      commit('SET_ORDER', order);
      try {
        await window.backend.Retro.SetGameOrder(game, order);
        console.log('Successfully set game order');
      } catch (e) {
        console.error('Failed to set game order:', e);
      }
    },
    async setGameMode({ commit }, hardcore) {
      commit('SET_MODE', hardcore);
      // TODO: send to backend
    },
  },
};
