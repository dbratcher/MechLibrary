function upvote(key) {
  console.log('upvote key', key);
  $.get('/mechanics/upvote?key='+key);

  // up vote count
  var count = $('#votes-'+key).html();
  $('#votes-'+key).html(parseInt(count, 10)+1);

  // remove vote link
  $('#upvoteLink-'+key).remove();

}
