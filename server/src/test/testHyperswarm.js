//Nimbus/server/src/test/testHyperswarm.js
const { expect } = require('chai');
const Hyperswarm = require('hyperswarm');
const { setupHyperswarm } = require('../api/hyperswarm');  // Import the module to be tested

describe('Hyperswarm Peer Discovery and Connections', function () {
  let swarm;

  before(function () {
    // Initialize the Hyperswarm instance
    swarm = Hyperswarm();
  });

  after(function () {
    swarm.destroy();  // Clean up the swarm after tests
  });

  it('should initialize Hyperswarm correctly', function () {
    expect(swarm).to.exist;
  });

  it('should discover peers successfully', function (done) {
    this.timeout(5000);  // Allow some time for peer discovery

    swarm.join(Buffer.from('mytopic', 'hex'), { announce: true });

    swarm.on('connection', (conn, info) => {
      expect(info).to.exist;
      expect(conn).to.exist;
      done();
    });
  });

  it('should handle peer disconnection', function (done) {
    swarm.on('disconnection', (conn) => {
      expect(conn).to.exist;
      done();
    });

    swarm.leave(Buffer.from('mytopic', 'hex'));
  });
});
