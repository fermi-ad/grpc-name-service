package nameserver.dbimport;

import java.io.BufferedReader;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class DbDirectoryParser {

    private static final Pattern RECORD_PATTERN =
            Pattern.compile("^\\s*record\\s*\\(\\s*([a-zA-Z0-9_]+)\\s*,\\s*\"([^\"]+)\"\\s*\\)");

    public List<PvRecord> parseDirectory(Path rootDir) throws IOException {
        List<Path> dbFiles = findDbFiles(rootDir);
        List<PvRecord> allRecords = new ArrayList<>();

        for (Path dbFile : dbFiles) {
            allRecords.addAll(parseFile(dbFile));
        }

        return allRecords;
    }

    private List<Path> findDbFiles(Path rootDir) throws IOException {
        try (Stream<Path> stream = Files.walk(rootDir)) {
            return stream
                    .filter(Files::isRegularFile)
                    .filter(path -> path.toString().endsWith(".db"))
                    .sorted()
                    .collect(Collectors.toList());
        }
    }

    private List<PvRecord> parseFile(Path dbFile) throws IOException {
        List<PvRecord> results = new ArrayList<>();

        try (BufferedReader reader = Files.newBufferedReader(dbFile)) {
            String line;
            int lineNumber = 0;

            while ((line = reader.readLine()) != null) {
                lineNumber++;

                Matcher matcher = RECORD_PATTERN.matcher(line);
                if (matcher.find()) {
                    String recordType = matcher.group(1);
                    String pvName = matcher.group(2);

                    results.add(new PvRecord(pvName, recordType, dbFile, lineNumber));
                }
            }
        }

        return results;
    }
}