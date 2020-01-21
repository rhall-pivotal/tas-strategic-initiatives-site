require("tap").mochaGlobals();
var chai = require('chai');
generateGuid = function() { return 'GUID'; };

chai.should();
chai.config.truncateThreshold = 0;
