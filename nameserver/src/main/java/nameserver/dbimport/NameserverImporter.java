package nameserver.dbimport;

import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;

public class NameserverImporter {

    public void importPvs(List<PvRecord> pvRecords) {
        Set<String> uniquePvs = new LinkedHashSet<>();

        for (PvRecord record : pvRecords) {
            uniquePvs.add(record.getPvName());
        }

        System.out.println("Preparing to import " + uniquePvs.size() + " unique PVs into Nameserver");

        for (String pvName : uniquePvs) {
            insertIntoNameserver(pvName);
        }

        System.out.println("Import finished");
    }

    private void insertIntoNameserver(String pvName) {
        // For now this is just a placeholder.
        // Later replace this with your real Nameserver service/repository/API call.
        System.out.println("INSERT INTO NAMESERVER: " + pvName);
    }
}