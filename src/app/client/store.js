import { compose, createStore as reduxCreateStore} from 'redux';
import { devTools, persistState } from 'redux-devtools';
import goStarterKit from './reducers';

let finalCreateStore;
if (process.env.NODE_ENV === 'production') {
  finalCreateStore = reduxCreateStore.bind(null, goStarterKit);
} else {
  try {
    finalCreateStore = compose(
      devTools(),
      persistState(window.location.href.match(/[?&]debug_session=([^&]+)\b/))
    )(reduxCreateStore).bind(null, goStarterKit);
    console.log('dev tools added');
  } catch (e) {
    finalCreateStore = compose(
      devTools()
    )(reduxCreateStore).bind(null, goStarterKit);
  }
}

export const createStore = finalCreateStore;
