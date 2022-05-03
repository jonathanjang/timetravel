import requests
import unittest
import pdb

class BasicTest( unittest.TestCase ):

    def verifyGetRecord( self, url, recordId, statusCode, expectedContent ):
        r = requests.get( url + str( recordId ) )
        self.assertEqual( r.status_code, statusCode )
        self.assertEqual( r.content, expectedContent )

    def verifyGetValuesRecord( self, url, recordId, key, statusCode, expectedContent ):
        r = requests.get( url + str( recordId ) + "/" + key )
        self.assertEqual( r.status_code, statusCode )
        self.assertEqual( r.content, expectedContent )

    def verifyPostRecord( self, url, payload, recordId, statusCode, expectedContent ):
        r = requests.post( url + str( recordId ), json=payload ) 
        self.assertEqual( r.status_code, statusCode )
        self.assertEqual( r.content, expectedContent )

    def testServer( self ):
        # Make sure to clean out .db file before running test!
        # Also make sure the server is running in the background "go run server.go"
        url = 'http://127.0.0.1:8000/api/v2/records/'

        # Verify that a GET with nothing in the database returns an error
        self.verifyGetRecord( url, 1, 400, '{"error":"record of id 1 does not exist"}\n' )

        # Test for add/update/delete for ID 1
        d = { "foo": "bar" }
        self.verifyPostRecord( url, d, 1, 200, '{"id":1,"data":{"foo":"bar"}}\n' )
    
        # follow every post with a GET to verify the GET response
        self.verifyGetRecord( url, 1, 200, '{"id":1,"data":{"foo":"bar"}}\n' )

        d = { "1234": "5678" }
        self.verifyPostRecord( url, d, 1, 200, '{"id":1,"data":{"1234":"5678","foo":"bar"}}\n' )
        self.verifyGetRecord( url, 1, 200, '{"id":1,"data":{"1234":"5678","foo":"bar"}}\n' )

        d = { "1234": None }
        self.verifyPostRecord( url, d, 1, 200, '{"id":1,"data":{"foo":"bar"}}\n' )
        self.verifyGetRecord( url, 1, 200, '{"id":1,"data":{"foo":"bar"}}\n' )

        # Test for add/update/delete for ID 2
        d = { "foo": "bar" }
        self.verifyPostRecord( url, d, 2, 200, '{"id":2,"data":{"foo":"bar"}}\n' )
        self.verifyGetRecord( url, 2, 200, '{"id":2,"data":{"foo":"bar"}}\n' )

        d = { "1234": "5678" }
        self.verifyPostRecord( url, d, 2, 200, '{"id":2,"data":{"1234":"5678","foo":"bar"}}\n' )
        self.verifyGetRecord( url, 2, 200, '{"id":2,"data":{"1234":"5678","foo":"bar"}}\n' )

        # Test case when a key is overwritten
        d = { "foo": "baz" }
        self.verifyPostRecord( url, d, 2, 200, '{"id":2,"data":{"1234":"5678","foo":"baz"}}\n' )
        self.verifyGetRecord( url, 2, 200, '{"id":2,"data":{"1234":"5678","foo":"baz"}}\n' )

        d = { "1234": None }
        self.verifyPostRecord( url, d, 2, 200, '{"id":2,"data":{"foo":"baz"}}\n' )
        self.verifyGetRecord( url, 2, 200, '{"id":2,"data":{"foo":"baz"}}\n' )

        d = { "foo": None }
        self.verifyPostRecord( url, d, 1, 200, '{"id":1,"data":{}}\n' )
        self.verifyGetRecord( url, 1, 200, '{"id":1,"data":{}}\n' )

        d = { "foo": None }
        self.verifyPostRecord( url, d, 2, 200, '{"id":2,"data":{}}\n' )
        self.verifyGetRecord( url, 2, 200, '{"id":2,"data":{}}\n' )

        # Test for GET on /records/{id}/{key} endpoint 
        # Add some more values to for a given key and verify functionality
        # of the new query on id, key
        d = { "foo": "34" }
        self.verifyPostRecord( url, d, 2, 200, '{"id":2,"data":{"foo":"34"}}\n' )
        self.verifyGetRecord( url, 2, 200, '{"id":2,"data":{"foo":"34"}}\n' )

        d = { "foo": "12" }
        self.verifyPostRecord( url, d, 2, 200, '{"id":2,"data":{"foo":"12"}}\n' )
        self.verifyGetRecord( url, 2, 200, '{"id":2,"data":{"foo":"12"}}\n' )

        d = { "foo": "78" }
        self.verifyPostRecord( url, d, 2, 200, '{"id":2,"data":{"foo":"78"}}\n' )
        self.verifyGetRecord( url, 2, 200, '{"id":2,"data":{"foo":"78"}}\n' )

        self.verifyGetValuesRecord( url, 2, "foo", 200, 
            '{"rid":2,"key":"foo","data":{"0":"78","1":"12","2":"34","3":"","4":"baz","5":"bar"}}\n' )

        
if __name__ == '__main__':
    unittest.main()
