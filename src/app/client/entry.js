const router = require('expose?__router__!./router');

if(typeof window !== 'undefined') {
  // Start main application here
  router.run();
}
