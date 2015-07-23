var express=require('express'),
  app=express(),
  https=require('https'),
  path=require('path'),
  swig=require('swig'),
  config=require('./config');

app.set('etag','strong');
app.set('port',(process.env.PORT||6000));
app.engine('svg',swig.renderFile);
app.set('views',path.join(__dirname,'/views'));
app.set('view cache',false);

app.get("/hackaday", function(req,res) {
  res.sendFile(path.join(__dirname,"/views/index.html"));
});

app.get('/hackaday/:id.svg', function(req,res) {
  var id=req.params.id.replace(/[|&;$%@"<>()+,]/g,"");
  var url="https://api.hackaday.io/v1/projects/"+id+"?api_key="+config.API_KEY;
  https.get(url, function(resp) {
    var body='';

    resp.on('data', function(chunk) {body+=chunk;});

    resp.on('end', function() {
      json_res=JSON.parse(body);
      if ('message' in json_res) {
        console.log('Package %s Gives Error %s',id,json_res['message']);
        res.status(500).send(json_res['message']);
      } else {
        res.type('image/svg+xml');
        res.append('Cache-Control','private, max-age=0, no-cache, no-store');
        res.append('Pragma','no-cache');
        res.render("base.svg",{json:json_res});
      }
    });
  }).on('error', function(err) {
    console.log('Error: %s',err);
    res.status(500).send(err);
  });
});

var server = app.listen(app.get('port'), function() {
  console.log("Node app is running at localhost:"+app.get('port'));
});

process.on('SIGINT', server.close);
process.on('SIGTERM', server.close);
