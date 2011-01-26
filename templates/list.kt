<a href='/edit/'>New article</a>
<hr>
<ul id='List'>
$for _, art in articles:
  <li><a href='$art.Data[id]'>$art.Data[title]</a></li>
$end
</ul>
