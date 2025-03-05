package main

var Commands = map[string]string{
	// Commands
	"ADMIN":    "ADMIN",    // Find admin info
	"AWAY":     "AWAY",     // Set/remove away message
	"CONNECT":  "CONNECT",  // Connect server to server
	"DIE":      "DIE",      // Shutdown server
	"ERROR":    "ERROR",    // Report error
	"INFO":     "INFO",     // Server information
	"INVITE":   "INVITE",   // Invite user to channel
	"ISON":     "ISON",     // Check if users are online
	"JOIN":     "JOIN",     // Join channel
	"KICK":     "KICK",     // Remove user from channel
	"KILL":     "KILL",     // Close client connection
	"LINKS":    "LINKS",    // List server connections
	"LIST":     "LIST",     // List channels and topics
	"LUSERS":   "LUSERS",   // Get user statistics
	"MODE":     "MODE",     // Change user/channel modes
	"MOTD":     "MOTD",     // Get Message of the Day
	"NAMES":    "NAMES",    // List users in channel
	"NICK":     "NICK",     // Set/change nickname
	"NOTICE":   "NOTICE",   // Send notice
	"OPER":     "OPER",     // Become an operator
	"PART":     "PART",     // Leave channel
	"PASS":     "PASS",     // Set connection password
	"PING":     "PING",     // Ping server/user
	"PONG":     "PONG",     // Pong reply
	"PRIVMSG":  "PRIVMSG",  // Send message
	"QUIT":     "QUIT",     // Disconnect
	"REHASH":   "REHASH",   // Reload server config
	"RESTART":  "RESTART",  // Restart server
	"SERVICE":  "SERVICE",  // Register service
	"SERVLIST": "SERVLIST", // List services
	"SQUERY":   "SQUERY",   // Service query
	"SQUIT":    "SQUIT",    // Disconnect server links
	"STATS":    "STATS",    // Query statistics
	"SUMMON":   "SUMMON",   // Summon user
	"TIME":     "TIME",     // Query local time
	"TOPIC":    "TOPIC",    // Channel topic
	"TRACE":    "TRACE",    // Trace route
	"USER":     "USER",     // Set user info
	"USERHOST": "USERHOST", // Get user info
	"USERS":    "USERS",    // List users
	"VERSION":  "VERSION",  // Get version
	"WALLOPS":  "WALLOPS",  // Send to ops
	"WHO":      "WHO",      // Query user info
	"WHOIS":    "WHOIS",    // Query user info
	"WHOWAS":   "WHOWAS",   // Query offline user

	// Numeric Replies
	"001": "RPL_WELCOME",         // Welcome to the network
	"002": "RPL_YOURHOST",        // Your host is
	"003": "RPL_CREATED",         // Server created on
	"004": "RPL_MYINFO",          // Server info
	"005": "RPL_ISUPPORT",        // Server supports
	"200": "RPL_TRACELINK",       // Link info
	"201": "RPL_TRACECONNECTING", // Try. to connect
	"202": "RPL_TRACEHANDSHAKE",  // Handshake info
	"203": "RPL_TRACEUNKNOWN",    // Unknown connection
	"204": "RPL_TRACEOPERATOR",   // Operator
	"205": "RPL_TRACEUSER",       // User
	"206": "RPL_TRACESERVER",     // Server
	"207": "RPL_TRACESERVICE",    // Service
	"208": "RPL_TRACENEWTYPE",    // New type
	"209": "RPL_TRACECLASS",      // Class
	"210": "RPL_TRACERECONNECT",  // Reconnect
	"211": "RPL_STATSLINKINFO",   // Link info
	"212": "RPL_STATSCOMMANDS",   // Commands
	"213": "RPL_STATSCLINE",      // C line
	"214": "RPL_STATSNLINE",      // N line
	"215": "RPL_STATSILINE",      // I line
	"216": "RPL_STATSKLINE",      // K line
	"217": "RPL_STATSQLINE",      // Q line
	"218": "RPL_STATSYLINE",      // Y line
	"219": "RPL_ENDOFSTATS",      // End of stats
	"221": "RPL_UMODEIS",         // User mode
	"231": "RPL_SERVICEINFO",     // Service info
	"232": "RPL_ENDOFSERVICES",   // End of services
	"233": "RPL_SERVICE",         // Service
	"234": "RPL_SERVLIST",        // Service list
	"235": "RPL_SERVLISTEND",     // Service list end
	"241": "RPL_STATSLLINE",      // L line
	"242": "RPL_STATSUPTIME",     // Server uptime
	"243": "RPL_STATSOLINE",      // O line
	"244": "RPL_STATSHLINE",      // H line
	"245": "RPL_STATSSLINE",      // S line
	"246": "RPL_STATSPING",       // Server ping
	"247": "RPL_STATSBLINE",      // B line
	"250": "RPL_STATSDLINE",      // D line
	"251": "RPL_LUSERCLIENT",     // Users
	"252": "RPL_LUSEROP",         // Operators
	"253": "RPL_LUSERUNKNOWN",    // Unknown connections
	"254": "RPL_LUSERCHANNELS",   // Channels
	"255": "RPL_LUSERME",         // Local users
	"256": "RPL_ADMINME",         // Admin info
	"257": "RPL_ADMINLOC1",       // Admin loc1
	"258": "RPL_ADMINLOC2",       // Admin loc2
	"259": "RPL_ADMINEMAIL",      // Admin email
	"261": "RPL_TRACELOG",        // Trace log
	"262": "RPL_TRACEEND",        // Trace end
	"263": "RPL_TRYAGAIN",        // Try again
	"265": "RPL_LOCALUSERS",      // Local users
	"266": "RPL_GLOBALUSERS",     // Global users
	"300": "RPL_NONE",            // None
	"301": "RPL_AWAY",            // Away
	"302": "RPL_USERHOST",        // Userhost
	"303": "RPL_ISON",            // ISON
	"305": "RPL_UNAWAY",          // No longer away
	"306": "RPL_NOWAWAY",         // Now away
	"311": "RPL_WHOISUSER",       // Whois user
	"312": "RPL_WHOISSERVER",     // Whois server
	"313": "RPL_WHOISOPERATOR",   // Whois operator
	"314": "RPL_WHOWASUSER",      // Whowas user
	"315": "RPL_ENDOFWHO",        // End of WHO
	"316": "RPL_WHOISCHANOP",     // Whois chanop
	"317": "RPL_WHOISIDLE",       // Whois idle
	"318": "RPL_ENDOFWHOIS",      // End of WHOIS
	"319": "RPL_WHOISCHANNELS",   // Whois channels
	"321": "RPL_LISTSTART",       // List start
	"322": "RPL_LIST",            // List
	"323": "RPL_LISTEND",         // List end
	"324": "RPL_CHANNELMODEIS",   // Channel mode
	"325": "RPL_UNIQOPIS",        // Unique op
	"331": "RPL_NOTOPIC",         // No topic
	"332": "RPL_TOPIC",           // Topic
	"333": "RPL_TOPICWHOTIME",    // Topic who time
	"341": "RPL_INVITING",        // Inviting
	"342": "RPL_SUMMONING",       // Summoning
	"346": "RPL_INVITELIST",      // Invite list
	"347": "RPL_ENDOFINVITELIST", // End invite list
	"348": "RPL_EXCEPTLIST",      // Exception list
	"349": "RPL_ENDOFEXCEPTLIST", // End exception
	"351": "RPL_VERSION",         // Version
	"352": "RPL_WHOREPLY",        // Who reply
	"353": "RPL_NAMREPLY",        // Names reply
	"361": "RPL_KILLDONE",        // Kill done
	"362": "RPL_CLOSING",         // Closing
	"363": "RPL_CLOSEEND",        // Close end
	"364": "RPL_LINKS",           // Links
	"365": "RPL_ENDOFLINKS",      // End of links
	"366": "RPL_ENDOFNAMES",      // End of names
	"367": "RPL_BANLIST",         // Ban list
	"368": "RPL_ENDOFBANLIST",    // End of ban list
	"369": "RPL_ENDOFWHOWAS",     // End of WHOWAS
	"371": "RPL_INFO",            // Info
	"372": "RPL_MOTD",            // MOTD
	"373": "RPL_INFOSTART",       // Info start
	"374": "RPL_ENDOFINFO",       // End of info
	"375": "RPL_MOTDSTART",       // MOTD start
	"376": "RPL_ENDOFMOTD",       // End of MOTD
	"381": "RPL_YOUREOPER",       // You are oper
	"382": "RPL_REHASHING",       // Rehashing
	"383": "RPL_YOURESERVICE",    // You are service
	"384": "RPL_MYPORTIS",        // My port is
	"391": "RPL_TIME",            // Time
	"392": "RPL_USERSSTART",      // Users start
	"393": "RPL_USERS",           // Users
	"394": "RPL_ENDOFUSERS",      // End of users
	"395": "RPL_NOUSERS",         // No users

	// Error Replies
	"401": "ERR_NOSUCHNICK",        // No such nick
	"402": "ERR_NOSUCHSERVER",      // No such server
	"403": "ERR_NOSUCHCHANNEL",     // No such channel
	"404": "ERR_CANNOTSENDTOCHAN",  // Cannot send
	"405": "ERR_TOOMANYCHANNELS",   // Too many channels
	"406": "ERR_WASNOSUCHNICK",     // Was no such nick
	"407": "ERR_TOOMANYTARGETS",    // Too many targets
	"408": "ERR_NOSUCHSERVICE",     // No such service
	"409": "ERR_NOORIGIN",          // No origin
	"411": "ERR_NORECIPIENT",       // No recipient
	"412": "ERR_NOTEXTTOSEND",      // No text to send
	"413": "ERR_NOTOPLEVEL",        // No toplevel
	"414": "ERR_WILDTOPLEVEL",      // Wildcard in toplevel
	"415": "ERR_BADMASK",           // Bad mask
	"421": "ERR_UNKNOWNCOMMAND",    // Unknown command
	"422": "ERR_NOMOTD",            // No MOTD
	"423": "ERR_NOADMININFO",       // No admin info
	"424": "ERR_FILEERROR",         // File error
	"431": "ERR_NONICKNAMEGIVEN",   // No nickname given
	"432": "ERR_ERRONEUSNICKNAME",  // Erroneous nickname
	"433": "ERR_NICKNAMEINUSE",     // Nickname in use
	"436": "ERR_NICKCOLLISION",     // Nick collision
	"437": "ERR_UNAVAILRESOURCE",   // Unavailable resource
	"441": "ERR_USERNOTINCHANNEL",  // User not in channel
	"442": "ERR_NOTONCHANNEL",      // Not on channel
	"443": "ERR_USERONCHANNEL",     // User on channel
	"444": "ERR_NOLOGIN",           // No login
	"445": "ERR_SUMMONDISABLED",    // Summon disabled
	"446": "ERR_USERSDISABLED",     // Users disabled
	"451": "ERR_NOTREGISTERED",     // Not registered
	"461": "ERR_NEEDMOREPARAMS",    // Need more params
	"462": "ERR_ALREADYREGISTRED",  // Already registered
	"463": "ERR_NOPERMFORHOST",     // No perm for host
	"464": "ERR_PASSWDMISMATCH",    // Password mismatch
	"465": "ERR_YOUREBANNEDCREEP",  // Banned
	"466": "ERR_YOUWILLBEBANNED",   // Will be banned
	"467": "ERR_KEYSET",            // Key already set
	"471": "ERR_CHANNELISFULL",     // Channel is full
	"472": "ERR_UNKNOWNMODE",       // Unknown mode
	"473": "ERR_INVITEONLYCHAN",    // Invite only
	"474": "ERR_BANNEDFROMCHAN",    // Banned from chan
	"475": "ERR_BADCHANNELKEY",     // Bad channel key
	"476": "ERR_BADCHANMASK",       // Bad channel mask
	"477": "ERR_NOCHANMODES",       // No channel modes
	"478": "ERR_BANLISTFULL",       // Ban list full
	"481": "ERR_NOPRIVILEGES",      // No privileges
	"482": "ERR_CHANOPRIVSNEEDED",  // Chan op needed
	"483": "ERR_CANTKILLSERVER",    // Can't kill server
	"484": "ERR_RESTRICTED",        // Restricted
	"485": "ERR_UNIQOPPRIVSNEEDED", // Uniq op needed
	"491": "ERR_NOOPERHOST",        // No oper host
	"501": "ERR_UMODEUNKNOWNFLAG",  // Mode unknown flag
	"502": "ERR_USERSDONTMATCH",    // Users don't match
}

var RCommands = map[string]string{}

func init() {
	for k, v := range Commands {
		RCommands[v] = k
	}
}
