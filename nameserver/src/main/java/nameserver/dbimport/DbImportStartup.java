package nameserver.dbimport;

import io.quarkus.runtime.StartupEvent;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.event.Observes;

@ApplicationScoped
public class DbImportStartup {

    void onStart(@Observes StartupEvent ev) {
        System.out.println("=== DbImportStartup reached ===");

        String enabled = System.getenv("IMPORT_DB");
        System.out.println("=== IMPORT_DB = " + enabled + " ===");

        if (!"true".equalsIgnoreCase(enabled)) {
            System.out.println("=== DB import is disabled ===");
            return;
        }

        try {
            DbImportRunner.runImport();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}