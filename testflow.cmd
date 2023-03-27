1) instantiate controller
2) using the controller instantiate mockstore
3) buildstab to tell the mock to controll the mocked server
3) using test pkg and mockstore instantiate server
4) instantiate recorder and request
5) call serveHttp of the server router
6) checkresponse