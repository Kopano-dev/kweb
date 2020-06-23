<?php
/**
 * Simple PHP script to simulate long running PHP scripts. Call with ?delay=300
 * and wait this long. The PHP request must complete after this delay and not
 * before. If it does complete earlier, then a timeout in the server config is
 * too low for the choosen duration.
 *
 * In kweb, it is sufficient to set read_timeout inside the corresponding
 * fastcgi or fastcgi2 block, which is used to serve such long running PHP
 * endpoints.
 *
 */

if (isset($_GET['delay'])) {
        $delay = intval($_GET['delay'], 10);
} else {
        $delay = 30;
}

set_time_limit(0);
sleep($delay);

?>
