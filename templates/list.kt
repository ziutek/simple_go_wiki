<a href='/edit/'>New article</a>
<hr>
<ul id='List'>
$for _, art in Articles:
  <li><a href='$art[Id]'>$art[Title]</a></li>
$end
</ul>
