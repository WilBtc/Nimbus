/test/testAtsign.js
const { expect } = require('chai');
const { authenticateAtsign, encryptMessage, decryptMessage } = require('../api/atsignAuth');  // Import functions to be tested

describe('Atsign Authentication and Encryption', function () {
  const atSign = '@example';
  const message = 'Hello, Nimbus!';
  let encryptedMessage;

  it('should authenticate with a valid atSign', async function () {
    const isAuthenticated = await authenticateAtsign(atSign);
    expect(isAuthenticated).to.be.true;
  });

  it('should encrypt the message correctly', function () {
    encryptedMessage = encryptMessage(atSign, message);
    expect(encryptedMessage).to.exist;
    expect(encryptedMessage).to.be.a('string');
  });

  it('should decrypt the message correctly', function () {
    const decryptedMessage = decryptMessage(atSign, encryptedMessage);
    expect(decryptedMessage).to.equal(message);
  });
});
