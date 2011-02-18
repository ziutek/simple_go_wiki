$if id:
  <div>
    <h2>$title</h2>
    $:markdown(body)
  </div>
  <div id='Actions'><a href='/edit/$id'>Edit</a></div>
$else:
  <h3>Simple Wiki</h3>
  <p>This application was wrote entirely in Go language, using the following
  external packages:</p>
  <ul>
    <li><a href='https://github.com/hoisie/web.go'>web.go</a></li>
    <li><a href='https://github.com/ziutek/kasia.go'>kasia.go</a></li>
    <li><a href='https://github.com/ziutek/kview'>kview</a></li>
    <li><a href='https://github.com/ziutek/mymysql'>MyMySQL</a></li>
  </ul>
$end
