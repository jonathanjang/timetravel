import requests
import unittest

class BasicTest( unittest.TestCase ):

    def testServer( self ):
        # Make sure to clean out .db file before running test!
        url = 'http://127.0.0.1:8000/api/v1/records/'

        r = requests.get( url + "1" )
        self.assertTrue( r.status_code, 400 )
        self.assertTrue( r.content, '{"error":"record of id 1 does not exist"}\n' )

        # FIXME: make this less redundant
        # Test for add/update/delete for ID 1
        d = { "foo": "bar" }
        r = requests.post( url + '1', json=d ) 
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":1,"data":{"foo":"bar"}}\n' )
    
        # Make sure GET still works
        r = requests.get( url + "1" )
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":1,"data":{"foo":"bar"}}\n' )

        d = { "1234": "5678" }
        r = requests.post( url + '1', json=d ) 
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":1,"data":{"1234":"5678","foo":"bar"}}\n' )

        r = requests.get( url + "1" )
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":1,"data":{"1234":"5678","foo":"bar"}}\n' )

        d = { "1234": None }
        r = requests.post( url + '1', json=d ) 
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":1,"data":{"foo":"bar"}}\n' )

        r = requests.get( url + "1" )
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":1,"data":{"foo":"bar"}}\n' )

        # Test for add/update/delete for ID 2
        d = { "foo": "bar" }
        r = requests.post( url + '2', json=d ) 
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":2,"data":{"foo":"bar"}}\n' )

        r = requests.get( url + "2" )
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":2,"data":{"foo":"bar"}}\n' )

        d = { "1234": "5678" }
        r = requests.post( url + '2', json=d ) 
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":2,"data":{"1234":"5678","foo":"bar"}}\n' )

        r = requests.get( url + "2" )
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":2,"data":{"1234":"5678","foo":"bar"}}\n' )

        d = { "1234": None }
        r = requests.post( url + '2', json=d ) 
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":2,"data":{"foo":"bar"}}\n' )

        r = requests.get( url + "2" )
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":2,"data":{"foo":"bar"}}\n' )

        d = { "foo": None }
        r = requests.post( url + '1', json=d ) 
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":1,"data":{}}\n' )

        r = requests.get( url + "1" )
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":1,"data":{}}\n' )

        d = { "foo": None }
        r = requests.post( url + '2', json=d ) 
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":2,"data":{}}\n' )

        r = requests.get( url + "2" )
        self.assertTrue( r.status_code, 200 )
        self.assertTrue( r.content, '{"id":2,"data":{}}\n' )




if __name__ == '__main__':
    unittest.main()
