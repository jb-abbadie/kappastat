$(function() {
    $('#div_list_streams').hide();
    $.getJSON("/api/following", listFollowing);
    console.log("initialized");

    $('#list_streams').change( function() {
        $('#curve_chart').fadeOut({'duration':500,'queue':true});
        var selected = $('#list_streams option:selected').attr("value");
        if ( selected !== undefined) {
        $.getJSON("/api/stat/" + selected , function(data) {

            var result = [];
            for(var i in data) {
                result.push([new Date(data[i]['Start']), data[i]['Viewer'], data[i]['Messages']]);
            }
            drawChart2(result);
        });
        $('#curve_chart').fadeIn({'duration':500,'queue':true});
        }
    });
});

function drawChart2(inp) {
    var options = {
        chart: {
            title: 'Test',
        },
        series : {
            0:{axis: 'Viewer'},
            1:{axis: 'Messages'}
        },
        axes: {
            y: {
                Viewer: {label: "Viewership"},
                Messages: {label: "Message"}
            }
        },
    };
    var data = new google.visualization.DataTable();
    data.addColumn('date', 'Time');
    data.addColumn('number', "Viewer");
    data.addColumn('number', "Chat Messages");
    data.addRows(inp);

    var chart = new google.charts.Line(document.getElementById('curve_chart'));
    chart.draw(data, google.charts.Line.convertOptions(options));
}
