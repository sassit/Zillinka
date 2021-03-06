var assert = require('assert');
const api = require("../restApiCall.js");

describe("Fetch rhine pegel level web through API in JSON format", function () {

  it("should get the Rhine pegel level yesterday at noon (12:00:00)", async function () {
    const d = new Date(Date.now() - 86400000);
    const yesterday = d.toISOString().substring(0,10); // yesterday in format yyyy-mm-dd
    let level = 0;
    level = await api.getRhineLevel(yesterday);
    console.log(`   .. level of Rhine on ${yesterday} at noon is: ${level}`);
    assert(level > 0, `level still 0 after call to api`);
  });
});
