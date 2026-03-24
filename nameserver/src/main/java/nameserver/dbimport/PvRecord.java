package nameserver.dbimport;

import java.nio.file.Path;

public class PvRecord {
    private final String pvName;
    private final String recordType;
    private final Path sourceFile;
    private final int lineNumber;

    public PvRecord(String pvName, String recordType, Path sourceFile, int lineNumber) {
        this.pvName = pvName;
        this.recordType = recordType;
        this.sourceFile = sourceFile;
        this.lineNumber = lineNumber;
    }

    public String getPvName() {
        return pvName;
    }

    public String getRecordType() {
        return recordType;
    }

    public Path getSourceFile() {
        return sourceFile;
    }

    public int getLineNumber() {
        return lineNumber;
    }

    @Override
    public String toString() {
        return "PvRecord{" +
                "pvName='" + pvName + '\'' +
                ", recordType='" + recordType + '\'' +
                ", sourceFile=" + sourceFile +
                ", lineNumber=" + lineNumber +
                '}';
    }
}