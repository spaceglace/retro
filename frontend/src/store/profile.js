export default {
  namespaced: true,

  state: {
    profiles: {},
    profile: {},
  },

  mutations: {
    SET_PROFILES(state, profiles) {
      state.profiles = profiles;
    },
    ADD_PROFILE(state, { name, key }) {
      state.profiles[name] = { name, key };
    },
    REMOVE_PROFILE(state, name) {
      delete state.profiles[name];
    },
    SET_ACTIVE_PROFILE(state, name) {
      state.profile = state.profiles[name];
    },
  },

  actions: {
    async getProfiles({ commit }) {
      try {
        const profiles = await window.backend.Config.GetProfiles();
        commit('SET_PROFILES', profiles);
      } catch (e) {
        console.error('Failed to fetch profiles:', e);
      }
    },
    async addProfile({ commit }, { name, key }) {
      // TODO: make sure profile doesn't exist yet
      commit('ADD_PROFILE', { name, key });
      commit('SET_ACTIVE_PROFILE', name);
      try {
        const result = await window.backend.Config.AddProfile(name, key);
        console.log('Successfully added profile:', result);
      } catch (e) {
        console.error('Failed to add profile:', e);
      }
    },
    async removeProfile({ commit, state, dispatch }, { name }) {
      commit('REMOVE_PROFILE', name);
      try {
        const result = await window.backend.Config.RemoveProfile(name);
        console.log('Successfully removed profile:', result);
      } catch (e) {
        console.error('Failed to remove profile:', e);
      }

      // if this was the active profile, switch to whatever's first in the list
      if (name === state.profile) {
        let newProfile = '';
        if (Object.keys(state.profiles).length > 0) {
          newProfile = Object.keys(state.profiles)[0];
        }
        dispatch('setActiveProfile', { name: newProfile });
      }
    },
    async getActiveProfile({ commit }) {
      try {
        const profile = await window.backend.Config.GetActiveProfile();
        console.log('got active profile:', profile);
        commit('SET_ACTIVE_PROFILE', profile);
      } catch (e) {
        console.error('Failed to get active profile:', e);
      }
    },
    async setActiveProfile({ commit }, { name }) {
      // TODO check if profile exists
      commit('SET_ACTIVE_PROFILE', name);
      try {
        const result = await window.backend.Config.SetActiveProfile(name);
        console.log('Successfully set active profile:', result);
      } catch (e) {
        console.error('Failed to set active profile:', e);
      }
    },
  },
};
