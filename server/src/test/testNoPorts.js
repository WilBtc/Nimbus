/test/testNoPorts.js
const { expect } = require('chai');
const { initializeTunnel, closeTunnel } = require('../api/noPorts');  // Import the module to be tested

describe('NoPorts Tunnel Initialization and Management', function () {
  let tunnel;

  it('should initialize a NoPorts tunnel successfully', async function () {
    tunnel = await initializeTunnel();
    expect(tunnel).to.exist;
    expect(tunnel.status).to.equal('connected');
  });

  it('should close the tunnel successfully', async function () {
    const result = await closeTunnel(tunnel);
    expect(result).to.be.true;
  });
});
