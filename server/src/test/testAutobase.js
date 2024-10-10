/test/testAutobase.js
const { expect } = require('chai');
const Autobase = require('autobase');
const { setupAutobase } = require('../api/autobase');  // Import the module to be tested

describe('Autobase Collaborative Data Sharing', function () {
  let autobase;

  before(function () {
    autobase = setupAutobase();
  });

  it('should initialize Autobase with multiple writers', function () {
    expect(autobase).to.exist;
    expect(autobase.writers).to.have.lengthOf(2);  // Assuming two writers in the setup
  });

  it('should write collaboratively to Autobase', async function () {
    const result = await autobase.append('Collaborative data');
    expect(result).to.exist;
  });

  it('should retrieve the collaborative data', async function () {
    const data = await autobase.get(0);
    expect(data).to.exist;
    expect(data.toString()).to.equal('Collaborative data');
  });
});
