package nameserver.dbimport;

import nameserver.ChannelController;

import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class DbImportRunner {

    private static final Pattern RECORD_PATTERN =
            Pattern.compile("^\\s*record\\s*\\(\\s*([a-zA-Z0-9_]+)\\s*,\\s*\"([^\"]+)\"\\s*\\)");

    public static void runImport() throws Exception {
        System.out.println("=== runImport() started ===");

        List<String> dbFiles = List.of(
                "epics-db/Glassman.db",
                "epics-db/IonSource.db"
        );

        Set<String> uniquePvs = new LinkedHashSet<>();

        for (String resourcePath : dbFiles) {
            System.out.println("=== Reading file: " + resourcePath + " ===");
            parseResourceFile(resourcePath, uniquePvs);
        }

        System.out.println("=== Total unique PVs found: " + uniquePvs.size() + " ===");

        for (String pvName : uniquePvs) {
            insertIntoNameserver(pvName);
        }

        System.out.println("=== Import finished ===");
    }

    private static void parseResourceFile(String resourcePath, Set<String> uniquePvs) throws Exception {
        InputStream inputStream =
                DbImportRunner.class.getClassLoader().getResourceAsStream(resourcePath);

        if (inputStream == null) {
            System.out.println("=== Resource not found: " + resourcePath + " ===");
            return;
        }

        try (BufferedReader reader = new BufferedReader(
                new InputStreamReader(inputStream, StandardCharsets.UTF_8))) {

            String line;
            int lineNumber = 0;
            int filePvCount = 0;

            while ((line = reader.readLine()) != null) {
                lineNumber++;

                Matcher matcher = RECORD_PATTERN.matcher(line);
                if (matcher.find()) {
                    String recordType = matcher.group(1);
                    String pvName = matcher.group(2);

                    System.out.println("Found PV in " + resourcePath +
                            " line " + lineNumber +
                            ": type=" + recordType +
                            ", pv=" + pvName); 

                    if (uniquePvs.add(pvName)) {
                        filePvCount++;
                    } else {
                        System.out.println("Skipping duplicate PV: " + pvName);
                    }
                }
            }

            System.out.println("=== Unique PVs from " + resourcePath + ": " + filePvCount + " ===");
        }
    }

    private static void insertIntoNameserver(String pvName) {
        String existingDeviceName = System.getenv("IMPORT_DEVICE_NAME");
        System.out.println("=== IMPORT_DEVICE_NAME = " + existingDeviceName + " ===");

        if (existingDeviceName == null || existingDeviceName.isBlank()) {
            throw new IllegalStateException("IMPORT_DEVICE_NAME is not set");
        }

        var channel = proto.Channel.newBuilder()
                .setName(pvName)
                .setDescription("Imported from db file")
                .setMetadata("")
                .setDeviceName(existingDeviceName)
                .build();

        var request = proto.CreateChannelRequest.newBuilder()
                .setChannel(channel)
                .build();

        var response = ChannelController.createChannel(request);

        System.out.println("Created channel: " + response.getChannel().getName());
    }
}

//duplicate handling added with Set<String>

