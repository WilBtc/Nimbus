/test/testHypercore.js
const { expect } = require('chai');
const Hypercore = require('hypercore');
const { setupHypercore } = require('../api/hypercore');  // Import the module to be tested

describe('Hypercore Data Storage and Immutability', function () {
  let core;

  before(function () {
    // Initialize the Hypercore
    core = setupHypercore('./test-storage');
  });

  it('should append data correctly', function (done) {
    core.append('Hello, Hypercore!', function (err) {
      expect(err).to.not.exist;
      done();
    });
  });

  it('should retrieve data immutably', function (done) {
    core.get(0, function (err, data) {
      expect(err).to.not.exist;
      expect(data.toString()).to.equal('Hello, Hypercore!');
      done();
    });
  });
});
